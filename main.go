package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"encoding/json"
	"os"
	"os/exec"
)

type Configs struct {
	aws_AccessKeyID string `env:"aws_access_key_id`
	aws_secretAccessKey  string `env:"aws_secret_access_key`
	iamRoleArn        string `env:"iam_role_arn`
	testingText string `env:"testingText"`
}

func getSecret(myRoleArn string, region string) {
	secretName := "KiuwanCredential"

	sess := session.Must(session.NewSession(&aws.Config{
                                          	Region: aws.String("us-east-1"),
                                          }))
	creds := stscreds.NewCredentials(sess, myRoleArn)
	svc := secretsmanager.New(sess, &aws.Config{Credentials: creds})

	req  := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	fmt.Println("Secret req ", req)

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(req)
fmt.Println("~~~~~~~~~~~~~~Testing1~~~~~~~~~~~")
fmt.Println("~~~~~~~~~~~~~~Testing1~~~~~~~~~~~result",err)
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
fmt.Println("~~~~~~~~~~~~~~Testing2~~~~~~~~~~~")
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
	}
fmt.Println("~~~~~~~~~~~~~~Testing2~~~~~~~~~~~secretString" , secretString )
fmt.Println("~~~~~~~~~~~~~~Testing2~~~~~~~~~~~decodedBinarySecret" , decodedBinarySecret )
var p map[string]interface{}
jsonData := []byte(secretString)
err = json.Unmarshal(jsonData, &p)
if err != nil{
			fmt.Println("Base64 Decode Error:", err)
			return
}
 fmt.Println(p)

 for key, value := range p {
     fmt.Printf("%s value is %v\n", key, value)
 }

var kuser, kpassword string
kuser = fmt.Sprint( p["KiuwanUser"])
kpassword = fmt.Sprint( p["KiuwanPassword"])

if err := EnvmanAdd("KIUWAN_USERNAME", kuser); err != nil {
		fmt.Println("Failed to store KIUWAN_USERNAME:", err)
		os.Exit(1)
	}

	if err := EnvmanAdd("KIUWAN_PASSWORD", kpassword); err != nil {
  		fmt.Println("Failed to store KIUWAN_USERNAME:", err)
  		os.Exit(1)
  	}
}

func EnvmanAdd(key, value string) error {
	args := []string{"add", "-k", key, "-v", value}
	return RunCommand("envman", args...)
}
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
  var cfg Configs

  if cfg.testingText != "" {
		  fmt.Println("'****************testing ':", cfg.testingText)
	}else
	{
	fmt.Println("****************")
	}
  fmt.Println("The following Inputs from bitrise.secrets")
  aws_AccessKeyID := cfg.aws_AccessKeyID

	fmt.Println("'aws_AccessKeyID':", aws_AccessKeyID)
 
  aws_secretAccessKey := cfg.aws_secretAccessKey
  fmt.Println("'aws_secretAccessKey':", aws_secretAccessKey)

  fmt.Println("AWS Secrets Manager")
	//iamRoleArn := cfg.iamRoleArn
	iamRoleArn :="arn:aws:iam::027962030681:role/miki_limitedaccess"
	fmt.Println("'iamRoleArn':", iamRoleArn)
	region := "us-east-1"
	getSecret(iamRoleArn, region)

}
