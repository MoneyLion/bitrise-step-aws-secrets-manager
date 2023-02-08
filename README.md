# Bitrise Step for AWS Secrets Manager

Bitrise Step to fetch secrets from AWS Secrets Manager.

View [changelog](./CHANGELOG.md).

## Usage

Include this Step in your workflow, for example:

```yaml
workflows:
  foo:
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
        - aws_access_key_id: $AWS_ACCESS_KEY_ID
        - aws_secret_access_key: $AWS_SECRET_ACCESS_KEY
        - aws_default_region: a-region-1
        - secret_list: |
            arn:aws:secret-1 # username # USERNAME
            arn:aws:secret-2 # password # PASSWORD
    - script@1:
        inputs:
        - content: |
            #!/bin/bash
            #
            # Access your secrets via $USERNAME and $PASSWORD
```

This fetches the secrets, and places the referenced values into the environment variables `USERNAME` and `PASSWORD`, which can then be used in the subsequent steps within the workflow.

### Step input

Specify the list of secrets to be fetched, under the `secret_list` input, with each secret value's key-value pair on its own line. The format to specify each pair is:

```
<Secret ARN> # <JSON object key> # <Environment variable>
```

If your SecretString is not JSON, the `<JSON object key>` can be omitted with `_`. See the PlainText SecretString example below for more details.

### JSON SecretString Example
For example, if the given the secret is JSON with an ARN `arn:aws:secret-1`, and a secret value:

```json
{
  "username": "admin",
  "password": "str0ngpassword"
}
```

Specifying this line in the secret list:

```
arn:aws:secret-1 # username # USERNAME
```

Fetches the secret, retrieves the JSON value under the key `username`, and stores that value in the `USERNAME` environment variable. `$USERNAME` will now contain the value `admin`.


### PlainText SecretString Example
If the given secret is plain text (not JSON) with an ARN `arn:aws:secret-1`, and the secret value:

```
SOME_SECRET_VALUE
```

Specifying this line in the secret list:

```
arn:aws:secret-1 # username # USERNAME
```

The key `username` is ignored and can be omitted with `_`. Fetches the secret, retrieves the value, and stores that value in the `USERNAME` environment variable. `$USERNAME` will now contain the value `SOME_SECRET_VALUE`.

### Authenticating with AWS

Supply AWS credentials and region configuration via the Step's input:

```yaml
workflows:
  foo:
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
        - aws_access_key_id: $AWS_ACCESS_KEY_ID
        - aws_secret_access_key: $AWS_SECRET_ACCESS_KEY
        - aws_default_region: a-region-1
        - secret_list: |
            ...
```

The credentials have to be stored in workflow secret.

You may also use an AWS named profile from shared configuration file, via `aws_profile` Step input:

```yaml
workflows:
  foo:
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
        - aws_profile: some-profile   # Like this
        - secret_list: |
            ...
```

To assume an IAM role before fetching secrets, you may specify the role's ARN via `aws_iam_role_arn` input:

```yaml
workflows:
  foo:
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
        - aws_access_key_id: $AWS_ACCESS_KEY_ID
        - aws_secret_access_key: $AWS_SECRET_ACCESS_KEY
        - aws_default_region: a-region-1
        - aws_iam_role_arn: 'arn:aws:role/some-role'  # Like this
        - secret_list: |
            secret-line-1
            secret-line-2
```

## Development note

Ensure the following is installed:

  - Go Programming Language
  - Bitrise CLI

### Setting up

  1. Clone this repository.

  1. Run `go mod vendor`.

  1. Create `.bitrise.secrets.yml` from the [sample](./.bitrise.secrets.sample.yml). Populate the necessary values.

  1. In [bitrise.yml](./bitrise.yml), under the step titled "Step Test", specify the list of secrets that you want to fetch. Update as well the subsequent script step that echoes the secrets, referencing the environment variables that you use.

  1. Run `bitrise run test` to test the Bitrise step.

### Publishing

  1. Bump the `BITRISE_STEP_VERSION` in [bitrise.yml](./bitrise.yml).

  1. Make a commit.

  1. Create an annotated Git tag.

  1. Push the commits and tags.

  1. Set `MY_STEPLIB_REPO_FORK_GIT_URL` in local file `bitrise.secrets.yml` to point to your forked StepLib repository.

  1. Run `bitrise run share-this-step`.
