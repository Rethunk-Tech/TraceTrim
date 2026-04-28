# Contributing to TraceTrim

Thank you for your interest in contributing to TraceTrim! This document outlines the process for contributing code, reporting bugs, and suggesting features.

## Development Workflow

### Prerequisites

- Go 1.26 or later
- Git

### Getting Started

1. **Fork and clone** the repository
2. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** and test thoroughly
4. **Run tests and lints** before committing:
   ```bash
   go test ./...
   go vet ./...
   ```

### Code Style

Follow Go conventions:
- Use `go fmt` for formatting
- Run `go vet` to catch common issues
- Write clear comments explaining WHY, not WHAT
- Use table-driven tests for multiple cases
- Handle errors explicitly (never ignore `err`)

### Testing

- **Unit tests** — Test individual components in isolation
- **Integration tests** — Test end-to-end with real data
- **Test coverage** — Aim for >80% on modified packages

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/parser
```

### Committing Changes

Use conventional commit messages:
- `feat:` for new features
- `fix:` for bug fixes
- `refactor:` for code reorganization
- `test:` for test additions/changes
- `docs:` for documentation
- `chore:` for maintenance tasks

Example:
```
feat: add Python traceback support

This adds detection and cleaning for Python tracebacks,
enabling the tool to work with Python error logs.
```

### Submitting Pull Requests

1. **Push your branch** to your fork
2. **Open a pull request** against `main`
3. **Describe your changes** — What problem does it solve?
4. **Reference any issues** — Link related GitHub issues
5. **Pass CI checks** — Ensure all tests and linters pass

## Reporting Bugs

Before reporting a bug:
1. Check if it's already reported
2. Run with `--verbose` to gather logs
3. Verify the issue on your platform

When reporting, include:
- Platform (Windows/macOS/Linux)
- Go version (`go version`)
- TraceTrim version
- Error message or unexpected behavior
- Steps to reproduce
- Example stack trace or config (if applicable)

## Suggesting Features

For new features:
1. Check existing issues for similar requests
2. Describe the use case — Why is this needed?
3. Provide examples if possible
4. Consider performance and security implications

## Areas for Enhancement

We welcome contributions in these areas:
- Additional stack trace formats (Python, Ruby, etc.)
- New configuration options
- Performance optimizations
- Documentation improvements
- Cross-platform compatibility fixes

## Security Issues

Do not open public issues for security vulnerabilities. Report to **security@rethunk.tech** with:
- Description of the vulnerability
- Affected component(s) and version(s)
- Steps to reproduce
- Potential impact
- Suggested fix (optional)

See [SECURITY.md](./SECURITY.md) for full details.

## Legal

By contributing to TraceTrim, you agree that your contributions will be licensed under the same MIT license as the project.

## Questions?

Open an issue on GitHub or contact the maintainers. We're here to help!

---

Happy contributing!
