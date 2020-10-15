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
	"time"
)

// Deletes an entire secret and all of its versions. You can optionally include a
// recovery window during which you can restore the secret. If you don't specify a
// recovery window value, the operation defaults to 30 days. Secrets Manager
// attaches a DeletionDate stamp to the secret that specifies the end of the
// recovery window. At the end of the recovery window, Secrets Manager deletes the
// secret permanently. At any time before recovery window ends, you can use
// RestoreSecret () to remove the DeletionDate and cancel the deletion of the
// secret. You cannot access the encrypted secret information in any secret that is
// scheduled for deletion. If you need to access that information, you must cancel
// the deletion with RestoreSecret () and then retrieve the information.
//
//     *
// There is no explicit operation to delete a version of a secret. Instead, remove
// all staging labels from the VersionStage field of a version. That marks the
// version as deprecated and allows Secrets Manager to delete it as needed.
// Versions that do not have any staging labels do not show up in
// ListSecretVersionIds () unless you specify IncludeDeprecated.
//
//     * The
// permanent secret deletion at the end of the waiting period is performed as a
// background task with low priority. There is no guarantee of a specific time
// after the recovery window for the actual delete operation to occur.
//
// Minimum
// permissions To run this command, you must have the following permissions:
//
//     *
// secretsmanager:DeleteSecret
//
// Related operations
//
//     * To create a secret, use
// CreateSecret ().
//
//     * To cancel deletion of a version of a secret before the
// recovery window has expired, use RestoreSecret ().
func (c *Client) DeleteSecret(ctx context.Context, params *DeleteSecretInput, optFns ...func(*Options)) (*DeleteSecretOutput, error) {
	stack := middleware.NewStack("DeleteSecret", smithyhttp.NewStackRequest)
	options := c.options.Copy()
	for _, fn := range optFns {
		fn(&options)
	}
	addawsAwsjson11_serdeOpDeleteSecretMiddlewares(stack)
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
	addOpDeleteSecretValidationMiddleware(stack)
	stack.Initialize.Add(newServiceMetadataMiddleware_opDeleteSecret(options.Region), middleware.Before)
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
			OperationName: "DeleteSecret",
			Err:           err,
		}
	}
	out := result.(*DeleteSecretOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DeleteSecretInput struct {

	// (Optional) Specifies that the secret is to be deleted without any recovery
	// window. You can't use both this parameter and the RecoveryWindowInDays parameter
	// in the same API call. An asynchronous background process performs the actual
	// deletion, so there can be a short delay before the operation completes. If you
	// write code to delete and then immediately recreate a secret with the same name,
	// ensure that your code includes appropriate back off and retry logic. Use this
	// parameter with caution. This parameter causes the operation to skip the normal
	// waiting period before the permanent deletion that AWS would normally impose with
	// the RecoveryWindowInDays parameter. If you delete a secret with the
	// ForceDeleteWithouRecovery parameter, then you have no opportunity to recover the
	// secret. It is permanently lost.
	ForceDeleteWithoutRecovery *bool

	// Specifies the secret that you want to delete. You can specify either the Amazon
	// Resource Name (ARN) or the friendly name of the secret. If you specify an ARN,
	// we generally recommend that you specify a complete ARN. You can specify a
	// partial ARN too—for example, if you don’t include the final hyphen and six
	// random characters that Secrets Manager adds at the end of the ARN when you
	// created the secret. A partial ARN match can work as long as it uniquely matches
	// only one secret. However, if your secret has a name that ends in a hyphen
	// followed by six characters (before Secrets Manager adds the hyphen and six
	// characters to the ARN) and you try to use that as a partial ARN, then those
	// characters cause Secrets Manager to assume that you’re specifying a complete
	// ARN. This confusion can cause unexpected results. To avoid this situation, we
	// recommend that you don’t create secret names ending with a hyphen followed by
	// six characters. If you specify an incomplete ARN without the random suffix, and
	// instead provide the 'friendly name', you must not include the random suffix. If
	// you do include the random suffix added by Secrets Manager, you receive either a
	// ResourceNotFoundException or an AccessDeniedException error, depending on your
	// permissions.
	//
	// This member is required.
	SecretId *string

	// (Optional) Specifies the number of days that Secrets Manager waits before it can
	// delete the secret. You can't use both this parameter and the
	// ForceDeleteWithoutRecovery parameter in the same API call. This value can range
	// from 7 to 30 days. The default value is 30.
	RecoveryWindowInDays *int64
}

type DeleteSecretOutput struct {

	// The ARN of the secret that is now scheduled for deletion.
	ARN *string

	// The friendly name of the secret that is now scheduled for deletion.
	Name *string

	// The date and time after which this secret can be deleted by Secrets Manager and
	// can no longer be restored. This value is the date and time of the delete request
	// plus the number of days specified in RecoveryWindowInDays.
	DeletionDate *time.Time

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func addawsAwsjson11_serdeOpDeleteSecretMiddlewares(stack *middleware.Stack) {
	stack.Serialize.Add(&awsAwsjson11_serializeOpDeleteSecret{}, middleware.After)
	stack.Deserialize.Add(&awsAwsjson11_deserializeOpDeleteSecret{}, middleware.After)
}

func newServiceMetadataMiddleware_opDeleteSecret(region string) awsmiddleware.RegisterServiceMetadata {
	return awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "secretsmanager",
		OperationName: "DeleteSecret",
	}
}
