# Security Policy

## Reporting Security Vulnerabilities

**DO NOT** open a public GitHub issue for security vulnerabilities. Instead, please report them responsibly to:

**Email:** security@rethunk.tech  
**Response SLA:** We aim to respond to security reports within 24 hours.

When reporting a vulnerability, please include:
- Description of the vulnerability
- Affected component(s) and version(s)
- Steps to reproduce (if applicable)
- Potential impact
- Suggested fix (optional)

## Supported Versions

TraceTrim is an active development project. Security updates are applied to:

| Version | Support Status | Update Cadence |
|---------|----------------|---|
| Latest | Active | Continuous |

Only the latest version receives security updates. Users are encouraged to upgrade to the latest release for security patches.

## Security Practices

### Input Validation

- **Clipboard content validation** — All clipboard content validated before processing
  - UTF-8 validation prevents binary data corruption
  - Size limits (default: 1MB) prevent memory exhaustion
  - Content sanitization removes potentially dangerous patterns
- **Stack trace patterns** — Strictly validated regex patterns
  - No dynamic regex construction
  - No evaluation of matched content
  - Safe frame deduplication logic

### Memory Safety

- **Go language** — Memory-safe language with garbage collection
- **Bounds checking** — Protobuf runtime prevents buffer overflows
- **Resource limits** — Configurable content size limits
- **Proper cleanup** — Defer statements ensure resource release

### Platform Security

- **Cross-platform library** — `golang.design/x/clipboard` handles platform differences
- **No privileged operations** — No root/admin access required
- **Clipboard permissions** — Only requires standard application permissions
- **No network access** — Application operates entirely locally

### Error Handling

- **Safe error reporting** — No sensitive data in error messages
- **Graceful degradation** — Errors don't expose system information
- **Content sanitization** — Errors don't include clipboard content snippets
- **Resource cleanup** — Errors trigger proper cleanup/shutdown

## Testing & Validation

- **Unit tests** — All modules have >80% test coverage
- **Integration tests** — End-to-end clipboard and parsing tests
- **Fuzzing** — Input validation tested against malformed clipboard content
- **Linting** — `go vet` and `golangci-lint` catch common mistakes
- **Vulnerability scanning** — `govulncheck` checks for known vulnerabilities

## Known Vulnerabilities

None currently known. Reports are welcome via security@rethunk.tech.

## Dependency Management

TraceTrim dependencies:
- **golang.design/x/clipboard** — Cross-platform clipboard library
- **Standard Go library** — Core functionality only

**Security checks:**
- `go mod verify` — Verify module checksums
- `go mod tidy` — Remove unused dependencies
- `govulncheck ./...` — Scan for known vulnerabilities
- **Dependabot** — Automated vulnerability alerts (if enabled)

## Threat Model

### Attack Vectors Considered

| Vector | Risk | Mitigation |
|--------|------|-----------|
| **Memory Exhaustion** | High | Content size limits, bounded processing, proper cleanup |
| **Pattern Injection** | Medium | Strict regex validation, no dynamic compilation |
| **Malformed UTF-8** | Medium | UTF-8 validation before processing |
| **Rapid Content Changes** | Low | Polling-based approach handles gracefully |
| **Resource Leaks** | Low | Defer statements, error handling |
| **Race Conditions** | Low | Single-threaded clipboard monitoring |

### Attack Vectors NOT Applicable

- **Network attacks** — No network access; local-only operation
- **File system attacks** — Config file only; no arbitrary file I/O
- **Privilege escalation** — Runs with user permissions only; no elevation
- **Code injection** — No dynamic code execution; regex only
- **Clipboard hijacking** — Application reads/writes only; no control over other apps
- **Key logging** — No keystroke monitoring; clipboard-only
- **GUI injection** — No UI rendering; terminal/pipe only

## Security Best Practices

### For Users

- **Permissions** — Grant clipboard access permissions normally for any clipboard app
- **Clipboard content** — Be aware clipboard content may be temporarily visible
- **Script mode** — Use `--script-mode` for processing sensitive traces in pipelines
- **Keep updated** — Upgrade to latest version for security patches

### For Developers

- **Regex safety** — Patterns are pre-compiled; no dynamic regex from user input
- **Size limits** — Always enforce configurable size limits on clipboard content
- **Error handling** — Errors are safe; no content leaked in messages
- **Testing** — Add tests for new stack trace formats including edge cases

## Content Handling

### Preserved

- ✅ Original error messages (unchanged)
- ✅ All unique stack frames (deduplicated only)
- ✅ Indentation and formatting (preserved exactly)
- ✅ Non-stack-trace content (unchanged)
- ✅ File paths and line numbers (unchanged)

### Removed

- ❌ Consecutive duplicate frames only (same function + file + line)
- ❌ Nothing else; all content is safe

### Never Executed

- Stack traces are parsed as plain text
- No code extraction or execution
- No evaluation of content
- No system commands from traces

## Configuration Security

### Safe Defaults

```yaml
clipboard-polling-interval: 500ms         # Reasonable polling rate
clipboard-max-content-size: 1048576       # 1MB limit
parser-min-stack-lines: 2                 # Avoid false positives
auto-detect-script-mode: true             # Detect non-interactive use
show-timestamp: true                      # Logging for audit trail
```

### Configuration File Security

- Config file is optional; defaults are safe
- No credentials in config (would be ignored)
- File permissions respected (config readable by user)
- Syntax errors trigger fallback to defaults

## Incident Response

In the event of a confirmed security vulnerability:
1. Impact assessment (severity, affected versions, scope)
2. Fix development (in private branch if critical)
3. Testing with regression tests
4. Security update release (version bump, changelogs)
5. User notification (security advisory if critical)
6. Post-incident review (prevent similar issues)

## Security Checklist

Before using TraceTrim in production:
- [ ] Understand clipboard content being processed
- [ ] Configure size limits appropriate for environment
- [ ] Test with sample stack traces first
- [ ] Review logs for any errors or anomalies
- [ ] Keep application updated to latest version
- [ ] Configure script mode if processing in pipelines

## Contact

- **Security Issues:** security@rethunk.tech
- **General Support:** support@rethunk.tech
- **Website:** https://rethunk.tech

---

**Last updated:** 2026-04-28
