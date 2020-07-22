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

func getSecret(myRoleArn string, region string) {
	secretName := "KiuwanCredential"
	//aws_AccessKeyID := "AKIAQNAVKXJMQJJ4I7W7"
	//aws_secretAccessKey := "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"
	//Create a Secrets Manager client
	/*
		svc := secretsmanager.New(session.New(&aws.Config{
			Credentials: credentials.NewStaticCredentials(
			aws_AccessKeyID,
			aws_secretAccessKey,
			"",),
		}),
		aws.NewConfig().WithRegion(region))
	*/
	/*sess := secretsmanager.New(session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
		aws_AccessKeyID,
		aws_secretAccessKey,
		"",),
	}),
	aws.NewConfig().WithRegion(region))
	*/
	// Create the credentials from AssumeRoleProvider to assume the role
	// referenced by the “myRoleARN” ARN. Prompt for MFA token from stdin.
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
	fmt.Println("The following Inputs from bitrise.secrets")
	aws_AccessKeyID := os.Getenv("aws_AccessKeyID")
	fmt.Println("'aws_AccessKeyID':", aws_AccessKeyID)
	aws_secretAccessKey := os.Getenv("aws_scretAccessKey")
	fmt.Println("'aws_secretAccessKey':", aws_secretAccessKey)

	fmt.Println("AWS Secrets Manager")

	iamRoleArn := "arn:aws:iam::027962030681:role/miki_limitedaccess"
	region := "us-east-2"
	getSecret(iamRoleArn, region)

	/*
		secretName := "Kiuwan-Jira-API"
		region := "us-east-1"
		secret, _ := GetSecret(secretName)
		fmt.Println("hello world")
		fmt.Println(secret)
	*/

	// Step 1

	//aws_secretAccessKey := "3Xu12+Bt7ewmeK7s3BNqDml1Zazcy3ebV3sUgvic"
	//iamRoleArn := "arn:aws:iam::027962030681:role/miki_limitedaccess"
	/*
		var secretName,secretVersion,region string
		flag.StringVar(&secretName, "name", "KiuwanCredential", "Name of secret")
		flag.StringVar(&secretVersion, "version", "AWSCURRENT", "Version Stage (default: AWSCURRENT)")
		flag.StringVar(&region, "region", "us-east-2", "AWS Region (default: us-east-1)")
		flag.Parse()
		if len(secretName) == 0 {
			fmt.Println("Error: Secret name required.")
			os.Exit(1)
		}
		sn := secretName
		fmt.Println(sn)

	*/

	/*
		sess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(sess, iamRoleArn)
		svc := secretsmanager.New(sess, &aws.Config{Credentials: creds})
		svc := secretsmanager.New(session.New(),
			aws.NewConfig().WithRegion(region))

	*/
	// First, access these input from Bitrise secrets:
	//	- AWS Access Key ID
	//	- AWS Secret Access Key
	//	- AWS IAM role ARN
	//
	// Using the input, and AWS SDK for Go, assume an IAM role.
	// This role assumption is required since the role has the
	// permission to read the secrest.

	// Step 2
	//
	// Accept the list of secrets to fetch from bitrise.yml.
	//
	// Using the input, and AWS SDK for Go, fetch the secrets from
	// AWS Secrets Manager.
	//
	// Store the results in a variable first.

	// Step 3
	//
	// Take the result from Step 2, and store it into environment variable.
	// This would allow the secrets to be used in subsequent build steps.
	// Might have to refer to Bitrise on how it propagates the variables
	// to the build steps.

	//os.Exit(0)
}
