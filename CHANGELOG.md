# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2021-07-16
### Added
- Support space around the `#` delimiter in each secret for readability.

### Changed
- Rename an environment variable that is used by the custom step, from `AWS_ROLE_ARN` to `SECRETS_AWS_ROLE_ARN`, to avoid usage clash with AWS SDK for Go.

## [0.1.0] - 2021-07-16
### Added
- Support more options in preparing AWS credentials, aside from using static credentials.

### Changed
- AWS IAM role ARN input is now optional.

## [0.0.2] - 2020-11-13
### Added
- Print logs when fetching and propagating secrets.
- Cache fetched secrets.

## [0.0.1] - 2020-10-15
### Added
- Core functionality to fetch secrets from AWS Secrets Manager.
