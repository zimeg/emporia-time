# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Repeat lookups until a sureness of at least 80.0% is reached
- Check the status of Emporia's API before performing command
- Gather API tokens from AWS Cognito using Emporia credentials
- Select a device from connected Emporia devices
- Helpful information is displayed on an empty input
- Tests for energy conversions and usage extrapolation
- Linting checks and automated tests added on GitHub Actions

### Changed

- Execution info is output to stderr
- Config stored in `~/.config/etime/settings.json`
- Config file respects `XDG_CONFIG_HOME` environment variable
- Repository name changed to `emporia-time`
- Go version updated to 1.19

### Removed

- Tokens stored in environment variables
- Instructions for manually configuring devices

## [0.0.1] - 2022-11-27

### Added

- Perform user commands or scripts (without interactivity)
- Output results and timing information from user command
- Display energy usage stats from command duration
- Setup developer scripts in Makefile
- Instruct setup with a README
- Added the MIT license
- Created a CHANGELOG

[Unreleased]: https://github.com/e-zim/emporia-time/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/e-zim/emporia-time/releases/tag/v0.0.1
