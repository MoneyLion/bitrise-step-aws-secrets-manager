# AWS Secrets Manager

Bitrise custom step to fetch secrets from AWS Secrets Manager.

## Development note

Setting up:

  1. Clone this repository.

  1. Run `go mod vendor`.

  1. Create `.bitrise.secrets.yml` from the [sample](./.bitrise.secrets.sample.yml). Populate the necessary values.

  1. Run `bitrise run test` to test the Bitrise step.
