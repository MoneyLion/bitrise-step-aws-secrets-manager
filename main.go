package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
	"os/exec"
)

const (
	AWS_ACCESS_KEY_ID     = "aws_access_key_id"
	AWS_SECRET_ACCESS_KEY = "aws_secret_access_key"
	AWS_DEFAULT_REGION    = "aws_default_region"
	AWS_ROLE_ARN          = "aws_role_arn"
)

type localConfig struct {
	awsAccessKeyId     string
	awsSecretAccessKey string
	awsDefaultRegion   string
	awsRoleArn         string
}

type simpleJson map[string]interface{}

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

func fetchSecrets(secretId string, awsConfig aws.Config) (secretString string, err error) {
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

func loadJson(stringData string) (result simpleJson, err error) {
	jsonData := []byte(stringData)
	err = json.Unmarshal(jsonData, &result)
	return
}

func exportEnvVar(data simpleJson, dataKey string, envVarKey string) (err error) {
	dataValue, ok := data[dataKey]
	if !ok {
		err = errors.New(dataKey + " not found in secret")
		return
	}
	c := exec.Command("envman", "add", "--key", envVarKey, "--value", dataValue.(string))
	err = c.Run()
	return
}

func main() {
	lcfg := localConfig{
		awsAccessKeyId:     os.Getenv(AWS_ACCESS_KEY_ID),
		awsSecretAccessKey: os.Getenv(AWS_SECRET_ACCESS_KEY),
		awsDefaultRegion:   os.Getenv(AWS_DEFAULT_REGION),
		awsRoleArn:         os.Getenv(AWS_ROLE_ARN),
	}

	var secretId = "--snip--"

	awsConfig, err := buildAWSConfig(lcfg)
	if err != nil {
		panic(err)
	}

	assumeRole(lcfg, &awsConfig)

	secretString, err := fetchSecrets(secretId, awsConfig)
	if err != nil {
		panic(err)
	}

	secretJson, err := loadJson(secretString)
	if err != nil {
		panic(err)
	}

	exportEnvVar(secretJson, "password", "TEST")
}
