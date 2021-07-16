# AWS Secrets Manager

Bitrise custom step to fetch secrets from AWS Secrets Manager.

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
