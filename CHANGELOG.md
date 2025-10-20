# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1] - 2025-10-20

### Added

- Initial release of TraceTrim
- Cross-platform clipboard monitoring using golang.design/x/clipboard (Windows, macOS, Linux)
- Automatic detection and cleaning of JavaScript/React stack traces
- Remove repetitive stack frames while preserving essential error information
- Script mode for use in shell pipelines
- Configuration file support with command-line flag overrides
- Comprehensive test coverage
- Enhanced version display showing tag, commit hash, and build type (e.g., v0.1-123456 for releases, v0.1-123456-dev for development builds)

### Fixed

- Clipboard access using cross-platform library
- Memory management and error handling
- Configuration validation and loading

### Changed

- Improved binary naming convention
- Enhanced release automation
- Better error handling and logging

## [Unreleased]
