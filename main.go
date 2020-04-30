package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2"
	"os"
)

func main() {
	fmt.Println("AWS Secrets Manager")

	// Step 1
	//
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

	os.Exit(0)
}
