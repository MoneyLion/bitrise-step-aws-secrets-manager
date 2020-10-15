// Code generated by smithy-go-codegen DO NOT EDIT.

package secretsmanager

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	smithy "github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	smithyhttp "github.com/awslabs/smithy-go/transport/http"
)

// Generates a random password of the specified complexity. This operation is
// intended for use in the Lambda rotation function. Per best practice, we
// recommend that you specify the maximum length and include every character type
// that the system you are generating a password for can support. Minimum
// permissions To run this command, you must have the following permissions:
//
//     *
// secretsmanager:GetRandomPassword
func (c *Client) GetRandomPassword(ctx context.Context, params *GetRandomPasswordInput, optFns ...func(*Options)) (*GetRandomPasswordOutput, error) {
	stack := middleware.NewStack("GetRandomPassword", smithyhttp.NewStackRequest)
	options := c.options.Copy()
	for _, fn := range optFns {
		fn(&options)
	}
	addawsAwsjson11_serdeOpGetRandomPasswordMiddlewares(stack)
	awsmiddleware.AddRequestInvocationIDMiddleware(stack)
	smithyhttp.AddContentLengthMiddleware(stack)
	AddResolveEndpointMiddleware(stack, options)
	v4.AddComputePayloadSHA256Middleware(stack)
	retry.AddRetryMiddlewares(stack, options)
	addHTTPSignerV4Middleware(stack, options)
	awsmiddleware.AddAttemptClockSkewMiddleware(stack)
	addClientUserAgent(stack)
	smithyhttp.AddErrorCloseResponseBodyMiddleware(stack)
	smithyhttp.AddCloseResponseBodyMiddleware(stack)
	stack.Initialize.Add(newServiceMetadataMiddleware_opGetRandomPassword(options.Region), middleware.Before)
	addRequestIDRetrieverMiddleware(stack)
	addResponseErrorMiddleware(stack)

	for _, fn := range options.APIOptions {
		if err := fn(stack); err != nil {
			return nil, err
		}
	}
	handler := middleware.DecorateHandler(smithyhttp.NewClientHandler(options.HTTPClient), stack)
	result, metadata, err := handler.Handle(ctx, params)
	if err != nil {
		return nil, &smithy.OperationError{
			ServiceID:     ServiceID,
			OperationName: "GetRandomPassword",
			Err:           err,
		}
	}
	out := result.(*GetRandomPasswordOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetRandomPasswordInput struct {

	// A string that includes characters that should not be included in the generated
	// password. The default is that all characters from the included sets can be used.
	ExcludeCharacters *string

	// Specifies that the generated password should not include digits. The default if
	// you do not include this switch parameter is that digits can be included.
	ExcludeNumbers *bool

	// Specifies that the generated password should not include lowercase letters. The
	// default if you do not include this switch parameter is that lowercase letters
	// can be included.
	ExcludeLowercase *bool

	// A boolean value that specifies whether the generated password must include at
	// least one of every allowed character type. The default value is True and the
	// operation requires at least one of every character type.
	RequireEachIncludedType *bool

	// Specifies that the generated password should not include uppercase letters. The
	// default if you do not include this switch parameter is that uppercase letters
	// can be included.
	ExcludeUppercase *bool

	// Specifies that the generated password should not include punctuation characters.
	// The default if you do not include this switch parameter is that punctuation
	// characters can be included. The following are the punctuation characters that
	// can be included in the generated password if you don't explicitly exclude them
	// with ExcludeCharacters or ExcludePunctuation: ! " # $ % & ' ( ) * + , - . / : ;
	// < = > ? @ [ \ ] ^ _ ` { | } ~
	ExcludePunctuation *bool

	// Specifies that the generated password can include the space character. The
	// default if you do not include this switch parameter is that the space character
	// is not included.
	IncludeSpace *bool

	// The desired length of the generated password. The default value if you do not
	// include this parameter is 32 characters.
	PasswordLength *int64
}

type GetRandomPasswordOutput struct {

	// A string with the generated password.
	RandomPassword *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addawsAwsjson11_serdeOpGetRandomPasswordMiddlewares(stack *middleware.Stack) {
	stack.Serialize.Add(&awsAwsjson11_serializeOpGetRandomPassword{}, middleware.After)
	stack.Deserialize.Add(&awsAwsjson11_deserializeOpGetRandomPassword{}, middleware.After)
}

func newServiceMetadataMiddleware_opGetRandomPassword(region string) awsmiddleware.RegisterServiceMetadata {
	return awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "secretsmanager",
		OperationName: "GetRandomPassword",
	}
}
