# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of TraceTrim
- Cross-platform clipboard monitoring (Windows, macOS, Linux)
- Automatic detection and cleaning of JavaScript/React stack traces
- Remove repetitive stack frames while preserving essential error information
- Comprehensive test coverage
- CI/CD pipeline with automated releases

### Fixed

- Platform-specific clipboard access implementations
- Memory management for Windows API integration
- Cross-compilation support for multiple architectures

### Changed

- Improved binary naming convention (`tracetrim-{platform}-{arch}`)
- Enhanced release automation with checksum verification
- Better error handling and logging
