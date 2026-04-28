# AGENTS.md — Developer onboarding for TraceTrim

TraceTrim is a cross-platform Go CLI tool that monitors the clipboard for JavaScript/React stack traces and automatically removes repetitive blocks, making error logs easier to read.

## Quick Navigation

- **For users:** See [README.md](./README.md) — installation, usage, configuration, troubleshooting
- **For developers:** This file — building, testing, extending
- **For security:** See [SECURITY.md](./SECURITY.md)

## Project Structure

```
cmd/tracetrim/                Entry point for CLI binary
internal/
  clipboard/                  Clipboard monitoring (cross-platform via golang.design/x/clipboard)
  parser/                     Stack trace parsing and cleaning
  config/                     Configuration management (YAML + flags)
  models/                     Shared data structures
tests/                        Integration tests and fixtures
Makefile                      Build targets
go.mod                        Module definition
```

## Development Workflow

### Build

```bash
make build       # Compile binary to ./bin/tracetrim
make test        # Run test suite
make fmt         # Format code
make lint        # Run linter
make vet         # Go vet analysis
```

### Running from Source

```bash
go run ./cmd/tracetrim -- --help
go run ./cmd/tracetrim -- --verbose

# Script mode (stdin/stdout)
echo "Error: Something
    at myFunc (file.js:10:5)
    at myFunc (file.js:10:5)" | go run ./cmd/tracetrim --
```

### Testing

Run tests with coverage:

```bash
make test
go test -cover ./...
go test -cover ./internal/parser
go test -cover ./internal/clipboard
```

Tests must pass before any PR merge. Aim for >80% coverage on modified packages.

## Key Modules

### Clipboard Monitor (`internal/clipboard/`)

Cross-platform clipboard monitoring using `golang.design/x/clipboard`:

- Polls clipboard at configurable interval (default: 500ms)
- Detects content changes
- Handles platform differences (Windows/macOS/Linux) transparently
- Size-limited (default: 1MB) to prevent memory exhaustion

When modifying:
- Update tests in `internal/clipboard/monitor_test.go`
- Test on multiple platforms if possible
- Consider polling interval trade-offs (responsiveness vs. CPU)

### Stack Trace Parser (`internal/parser/`)

Detects and cleans JavaScript/React stack traces:

- Pattern matching for error types (Error, TypeError, ReferenceError, etc.)
- Frame deduplication (removes consecutive identical frames)
- Preserves formatting and non-stack-trace content
- Adds repeat count annotation (`// [x4]`) for cleaned duplicates

When modifying:
- Add test fixtures in `internal/parser/testdata/`
- Test regex patterns with edge cases
- Verify formatting is preserved
- Document new stack trace formats supported

### Configuration (`internal/config/`)

Loads configuration from file (`config.yaml`) and command-line flags:

- Flag precedence: CLI flags > config file > defaults
- Validates configuration values
- Supports auto-detection of script mode (TTY detection)

When modifying:
- Update config schema in `internal/config/config.go`
- Test YAML parsing with `config_test.go`
- Document new options in README.md

## Adding Features

### New Stack Trace Format

To support a new stack trace format (e.g., Python tracebacks):

1. Add pattern to `internal/parser/patterns.go`
2. Implement frame parser in `internal/parser/parser.go`
3. Add test fixtures in `internal/parser/testdata/`
4. Test with `go test -v ./internal/parser`
5. Update README § Supported Stack Trace Formats
6. Commit: `feat: add Python traceback support`

### New Configuration Option

To add a new option (e.g., `--max-repeat-count`):

1. Add field to config struct in `internal/config/config.go`
2. Add flag parsing with `flag.IntVar()`
3. Add YAML parsing support
4. Use option in appropriate module
5. Add tests covering the new option
6. Update README § Configuration
7. Commit: `feat: add max-repeat-count option`

### Clipboard Polling Optimization

To improve clipboard polling efficiency:

1. Consider debouncing logic (avoid re-processing same content)
2. Update interval with `--clipboard-polling-interval` flag
3. Profile performance with `go test -bench ./...`
4. Test on resource-constrained systems
5. Document trade-offs in README

## Dependency Management

Dependencies are minimal by design:

- **golang.design/x/clipboard** — Cross-platform clipboard access (required)
- **Standard library only** — Core parsing and config via YAML + flags

Installation:

```bash
go mod download
go mod verify
```

Before adding a dependency:
1. Justify in the commit message (why needed, why not built-in)
2. Check for security issues: `go mod tidy && govulncheck ./...`
3. Verify it doesn't bloat the binary (aim for <10MB)
4. Test on all supported platforms

## Testing Strategy

### Unit Tests

Test individual components in isolation:

```bash
go test -v ./internal/parser     # Parser tests
go test -v ./internal/config     # Config loading tests
go test -v ./internal/clipboard  # Clipboard tests
```

### Integration Tests

Test end-to-end with real clipboard:

```bash
go test -v ./tests/integration/
```

### Fixtures

Stack trace samples in `internal/parser/testdata/`:
- `javascript-simple.txt` — Basic error with duplicate frames
- `react-dev.txt` — React development error
- `node-native.txt` — Node.js native error
- etc.

Add new fixtures for new formats:
1. Create `internal/parser/testdata/<format>.txt`
2. Add to test file: `TestParser/formats`
3. Include expected output after cleaning

## Benchmarking

For performance-critical paths (parsing, clipboard polling):

```bash
# Run benchmarks
go test -bench ./internal/parser -benchmem

# Compare before/after:
go test -bench ./... -benchout=before.txt
# Make changes
go test -bench ./... -benchout=after.txt
go test -benchstat before.txt after.txt
```

Target: parse large stack traces (<50ms).

## Cross-Platform Considerations

TraceTrim runs on Windows, macOS, and Linux:

- **Clipboard access** — Handled by `golang.design/x/clipboard`
- **Path handling** — Use `filepath` package (not hardcoded `/` or `\`)
- **Line endings** — Normalize to `\n` internally; preserve in output
- **Permissions** — No special permissions needed; clipboard access only

Test on each platform when:
- Modifying clipboard monitor
- Adding new config options
- Changing parsing logic

## Code Style

- **Go conventions:** Follow `go fmt` and `go vet`
- **Naming:** Clear, short names (prefer `ctx` for context, `err` for error)
- **Comments:** Explain WHY, not WHAT (code is self-documenting)
- **Error handling:** Always check errors, use `fmt.Errorf()` for context
- **Tests:** Separate concerns, use table-driven tests for multiple cases

## Performance Targets

- **Binary size** — <10MB (including dependencies)
- **Memory usage** — <50MB at rest
- **Clipboard polling** — <2% CPU at 500ms interval
- **Stack trace parsing** — <50ms for 100-frame traces
- **Startup time** — <100ms

## Debugging

Enable verbose output:

```bash
go run ./cmd/tracetrim -- --verbose
```

Common issues:
- **Clipboard not accessible** — Platform-specific permissions
- **Patterns not matching** — Regex edge cases; update testdata
- **High CPU usage** — Increase polling interval or check for infinite loops

## References

- **[golang.design/x/clipboard](https://github.com/golang-design/clipboard)** — Clipboard library docs
- **[Go regexp package](https://golang.org/pkg/regexp/)** — Pattern matching reference
- **[Stack Trace Formats](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error/stack)** — MDN error documentation

---

**Last updated:** 2026-04-28
