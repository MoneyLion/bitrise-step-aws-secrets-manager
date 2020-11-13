package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"os/exec"
	"strings"
)

const (
	AWS_ACCESS_KEY_ID     = "aws_access_key_id"
	AWS_SECRET_ACCESS_KEY = "aws_secret_access_key"
	AWS_DEFAULT_REGION    = "aws_default_region"
	AWS_ROLE_ARN          = "aws_role_arn"
	SECRET_LIST           = "secret_list"
)

type localConfig struct {
	awsAccessKeyId     string
	awsSecretAccessKey string
	awsDefaultRegion   string
	awsRoleArn         string
	secretList         string
}

type secretListItem struct {
	arn    string
	key    string
	envvar string
}

type secretValueJson map[string]string

type secretCacheMap map[string]string

func buildAWSConfig(lcfg localConfig) (awsConfig aws.Config, err error) {
	awsConfig, err = config.LoadDefaultConfig(
		config.WithRegion(lcfg.awsDefaultRegion),
		config.WithCredentialsProvider{
			CredentialsProvider: credentials.NewStaticCredentialsProvider(
				lcfg.awsAccessKeyId,
				lcfg.awsSecretAccessKey,
				"",
			),
		},
	)
	return
}

func assumeRole(lcfg localConfig, awsConfig *aws.Config) {
	stsSvc := sts.NewFromConfig(*awsConfig)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, lcfg.awsRoleArn)
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
			arn:    secretComponents[0],
			key:    secretComponents[1],
			envvar: secretComponents[2],
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

func exportEnvVar(data secretValueJson, dataKey string, envVarKey string) (err error) {
	dataValue, ok := data[dataKey]
	if !ok {
		err = errors.New(dataKey + " not found in secret")
		return
	}
	fmt.Printf("Storing secret value for key '%s' into $%s\n", dataKey, envVarKey)
	c := exec.Command("envman", "add", "--key", envVarKey, "--value", dataValue)
	err = c.Run()
	return
}

func main() {
	lcfg := localConfig{
		awsAccessKeyId:     os.Getenv(AWS_ACCESS_KEY_ID),
		awsSecretAccessKey: os.Getenv(AWS_SECRET_ACCESS_KEY),
		awsDefaultRegion:   os.Getenv(AWS_DEFAULT_REGION),
		awsRoleArn:         os.Getenv(AWS_ROLE_ARN),
		secretList:         os.Getenv(SECRET_LIST),
	}

	awsConfig, err := buildAWSConfig(lcfg)
	if err != nil {
		panic(err)
	}

	assumeRole(lcfg, &awsConfig)

	secretCache := make(secretCacheMap)

	for _, item := range parseSecretList(lcfg.secretList) {
		secretString, err := cacher(secretCache, item.arn, awsConfig, fetchSecrets)
		if err != nil {
			panic(err)
		}

		secretJson, err := loadJson(secretString)
		if err != nil {
			panic(err)
		}

		exportEnvVar(secretJson, item.key, item.envvar)
	}
}
