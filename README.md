# AWS Secrets Manager

Bitrise custom step to fetch secrets from AWS Secrets Manager.

View [changelog](./CHANGELOG.md).

## Usage

Include this step in your workflow, for example:

```yaml
workflows:
  foo:
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
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

For example, given the secret with an ARN `arn:aws:secret-1`, and a secret value:

```
{
  "username": "admin",
  "password": "str0ngpassword"
}
```

Specifying this line in the secret list:

```
arn:aws:secret-1 # username # USERNAME
```

Fetches the secret, retrieves the JSON value under the key `username`, and store that value in the `USERNAME` environment variable. `$USERNAME` will now contain the value `admin`.

### Authenticating with AWS

The custom step uses AWS SDK for Go v2 with the default config loader. This means for authenticating with AWS, you may:

  - Use static AWS credentials via environment variable, e.g. `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.
  - Use shared configuration files, e.g. `AWS_PROFILE`.

To assume an IAM role before fetching secrets, you may specify the role's ARN via `AWS_IAM_ROLE_ARN` environment variable, or use the `aws_iam_role_arn` input:

```yaml
workflows:
  foo:
    envs:
      - AWS_IAM_ROLE_ARN: "arn:aws:role/some-role"  # This works
    steps:
    - aws-secrets-manager@x.x.x:
        inputs:
        - aws_iam_role_arn: "arn:aws:role/some-role"  # This works too
        - secret_list: |
            secret-line-1
            secret-line-2
```

## Development note

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
