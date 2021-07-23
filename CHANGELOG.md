# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Document on automatic sourcing of AWS configurations on Step's input

## [1.0.0] - 2021-07-22
### Added
- Support more options in preparing AWS credentials, aside from using static credentials.
- Support space around the `#` delimiter in each secret for readability.

### Changed
- Step input `aws_role_arn` is now optional.
- Step input `aws_role_arn` is renamed to `aws_iam_role_arn`.

### Removed
- No longer treat the value of `AWS_ROLE_ARN` environment variable as Step's input value.

## [0.0.2] - 2020-11-13
### Added
- Print logs when fetching and propagating secrets.
- Cache fetched secrets.

## [0.0.1] - 2020-10-15
### Added
- Core functionality to fetch secrets from AWS Secrets Manager.
