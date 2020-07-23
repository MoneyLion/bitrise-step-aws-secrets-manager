package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"os"
)

type config struct {
	aws_AccessKeyID string `env:"access_key_id"`
	aws_secretAccessKey  string `env:"secret_access_key"`
	iamRoleArn        string `env:"iam_role_arn"`
	testconfig  string `mikitest`
}

func getSecret(myRoleArn string, region string) {
	secretName := "KiuwanCredential"
	sess := session.Must(session.NewSession())
	creds := stscreds.NewCredentials(sess, myRoleArn)

	svc := secretsmanager.New(sess, &aws.Config{Credentials: creds})

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	fmt.Println("Secret ID is", input.SecretId)
	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
		fmt.Println("Secret string ", secretString)
		fmt.Println("Dedecoded Binary Secret", decodedBinarySecret)
	}
}

func main() {

	var cfg config
  	fmt.Println("The following Inputs from bitrise.secrets")
  	aws_AccessKeyID := cfg.aws_AccessKeyID
  	fmt.Println("'aws_AccessKeyID':", aws_AccessKeyID)
  	aws_secretAccessKey := cfg.aws_secretAccessKey
  	fmt.Println("'aws_secretAccessKey':", aws_secretAccessKey)
  	fmt.Println("************* test config : ", cfg.testconfig)

  	fmt.Println("AWS Secrets Manager")


  	iamRoleArn := cfg.iamRoleArn
  	ß
	region := "us-east-2"

	getSecret(iamRoleArn, region)

}
