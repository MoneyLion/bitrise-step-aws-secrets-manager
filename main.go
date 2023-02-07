package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type stepInput struct {
	awsAccessKeyId     string
	awsSecretAccessKey string
	awsDefaultRegion   string
	awsProfile         string
	awsIamRoleArn      string
	secretList         string
}

type secretListItem struct {
	arn    string
	key    string
	envvar string
}

type secretValueJson map[string]string

type secretCacheMap map[string]string

func prepareAwsConfig(sinput stepInput) (awsConfig aws.Config, err error) {
	if sinput.awsAccessKeyId != "" && sinput.awsSecretAccessKey != "" && sinput.awsDefaultRegion != "" {
		fmt.Println("Loading AWS config using static credentials")
		awsConfig, err = config.LoadDefaultConfig(
			config.WithRegion(sinput.awsDefaultRegion),
			config.WithCredentialsProvider{
				CredentialsProvider: credentials.NewStaticCredentialsProvider(
					sinput.awsAccessKeyId,
					sinput.awsSecretAccessKey,
					"",
				),
			})
	} else if sinput.awsProfile != "" {
		fmt.Println("Loading AWS config using named profile")
		awsConfig, err = config.LoadDefaultConfig(
			config.WithSharedConfigProfile(sinput.awsProfile))
	} else {
		err = errors.New("Incomplete AWS configuration. Specify AWS static credentials and region, or an AWS named profile for shared configuration, via the Step's input.")
	}

	return
}

func assumeRole(sinput stepInput, awsConfig *aws.Config) {
	stsSvc := sts.NewFromConfig(*awsConfig)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, sinput.awsIamRoleArn)
	awsConfig.Credentials = &aws.CredentialsCache{Provider: creds}
	return
}

func parseSecretList(secretList string) (items []secretListItem) {
	for _, secret := range strings.Split(secretList, "\n") {
		if strings.TrimSpace(secret) == "" {
			continue
		}
		secretComponents := strings.Split(secret, "#")
		items = append(items, secretListItem{
			arn:    strings.TrimSpace(secretComponents[0]),
			key:    strings.TrimSpace(secretComponents[1]),
			envvar: strings.TrimSpace(secretComponents[2]),
		})
	}
	return
}

func cacher(secretCache secretCacheMap, secretId string, awsConfig aws.Config, fetcher func(string, aws.Config) (string, error)) (secretString string, err error) {
	secretString, cached := secretCache[secretId]
	if cached {
		err = nil
		return
	}
	secretString, err = fetcher(secretId, awsConfig)
	if err != nil {
		return
	}
	secretCache[secretId] = secretString
	return
}

func fetchSecrets(secretId string, awsConfig aws.Config) (secretString string, err error) {
	fmt.Printf("Getting secret for %s\n", secretId)
	secretString = ""
	smSvc := secretsmanager.NewFromConfig(awsConfig)
	secretValue, err := smSvc.GetSecretValue(
		context.Background(),
		&secretsmanager.GetSecretValueInput{
			SecretId: &secretId,
		},
	)
	if err != nil {
		return
	}
	if secretValue.SecretString == nil {
		err = errors.New("Missing SecretString on GetSecretValue response")
		return
	}
	secretString = *secretValue.SecretString
	return
}

func loadJson(secretString string) (result secretValueJson, err error) {
	jsonData := []byte(secretString)
	err = json.Unmarshal(jsonData, &result)
	return
}

func exportEnvVarJson(data secretValueJson, dataKey string, envVarKey string) (err error) {
	dataValue, ok := data[dataKey]
	if !ok {
		err = errors.New(dataKey + " not found in secret")
		return
	}
	fmt.Printf("Storing secret value for key '%s' into $%s\n", dataKey, envVarKey)
	err = exportEnvVar(dataValue, envVarKey)
	return
}

func exportEnvVar(value string, envVarKey string) (err error) {
	c := exec.Command("envman", "add", "--key", envVarKey, "--value", value, "--sensitive")
	err = c.Run()
	return
}

func IsJSON(str string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(str), &js) == nil
}

func main() {
	// Prevent environment variables from interfering with AWS config loading
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_DEFAULT_REGION")
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_PROFILE")

	sinput := stepInput{
		awsAccessKeyId:     os.Getenv("aws_access_key_id"),
		awsSecretAccessKey: os.Getenv("aws_secret_access_key"),
		awsDefaultRegion:   os.Getenv("aws_default_region"),
		awsProfile:         os.Getenv("aws_profile"),
		awsIamRoleArn:      os.Getenv("aws_iam_role_arn"),
		secretList:         os.Getenv("secret_list"),
	}

	awsConfig, err := prepareAwsConfig(sinput)
	if err != nil {
		panic(err)
	}

	if sinput.awsIamRoleArn != "" {
		assumeRole(sinput, &awsConfig)
	}

	secretCache := make(secretCacheMap)

	for _, item := range parseSecretList(sinput.secretList) {
		secretString, err := cacher(secretCache, item.arn, awsConfig, fetchSecrets)
		if err != nil {
			panic(err)
		}

		if IsJSON(secretString) {
			secretJson, err := loadJson(secretString)
			if err != nil {
				panic(err)
			}

			exportEnvVarJson(secretJson, item.key, item.envvar)
		} else {
			exportEnvVar(secretString, item.envvar)
		}
	}
}
