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

## [v0.2.0] - 2026-06-11

### Added

- Makefile with cross-compile, test, and checksum targets; CI workflows updated for better portability

### Fixed

- golangci-lint v2 migration: config format, `goimports` local-prefixes, `noctx` and `goimports` linter errors in cmd package
- Reduced cyclomatic complexity of `main()` to satisfy lint v2 threshold
- Missing golangci-lint installation step in CI lint job

### Changed

- Go toolchain upgraded to 1.26; module dependencies refreshed across the board
- GitHub Actions pins updated to current majors (cache v5, gh-release v3, setup-go v6, checkout v6)
- Added HUMANS.md, AGENTS.md, CONTRIBUTING.md, and condensed README to governance pattern
- Security policy and CLAUDE.md symlink added for developer onboarding
- CI lint job migrated to golangci-lint v2 module path

> **Tag format normalized:** prior tag `v0.1` was two-part; the release workflow requires `vX.Y.Z`. This release adopts `v0.2.0` going forward.

## [Unreleased]
