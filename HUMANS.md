# HUMANS.md — User guide for TraceTrim

This file covers **installation, usage, and configuration**. For other documentation:
- **For developers:** See [AGENTS.md](./AGENTS.md) — architecture, testing, extending
- **For security:** See [SECURITY.md](./SECURITY.md) — vulnerability reporting, threat model
- **For contributors:** See [CONTRIBUTING.md](./CONTRIBUTING.md) — development workflow
- **Overview:** See [README.md](./README.md)

## Installation

### From Releases

Download pre-built binary from [Releases](https://github.com/Rethunk-Tech/TraceTrim/releases):

```bash
# Linux x86_64
wget https://github.com/Rethunk-Tech/TraceTrim/releases/download/v1.0.0/tracetrim-linux-amd64
chmod +x tracetrim-linux-amd64
sudo mv tracetrim-linux-amd64 /usr/local/bin/tracetrim

# macOS (Apple Silicon)
wget https://github.com/Rethunk-Tech/TraceTrim/releases/download/v1.0.0/tracetrim-darwin-arm64
chmod +x tracetrim-darwin-arm64
sudo mv tracetrim-darwin-arm64 /usr/local/bin/tracetrim

# Windows
# Download tracetrim-windows-amd64.exe and add to PATH
```

### Build from Source

```bash
git clone https://github.com/Rethunk-Tech/TraceTrim.git
cd TraceTrim
go build -o tracetrim ./cmd/tracetrim
sudo cp tracetrim /usr/local/bin/
```

## Usage

### Interactive Mode (Default)

Run TraceTrim to monitor clipboard:

```bash
tracetrim
```

The tool monitors your clipboard continuously. When you copy a JavaScript/React stack trace, it automatically:
1. Detects the trace
2. Removes repetitive frames
3. Updates clipboard with cleaned trace
4. Shows notification

**Example output:**

```
TraceTrim
Monitoring clipboard for JavaScript/React stack traces...
Press Ctrl+C to exit

[14:23:45] Detected stack trace, cleaning...
✓ Stack trace cleaned and clipboard updated
  Removed 3 repetitive lines
```

### Script Mode

Use TraceTrim in pipelines (auto-detected or explicit):

```bash
# From file
cat stack_trace.txt | tracetrim > cleaned.txt

# From pipe
echo "Error: Something
    at myFunc (file.js:10:5)
    at myFunc (file.js:10:5)" | tracetrim

# Explicit script mode
tracetrim --script-mode < stack.txt > clean.txt
```

### Command-Line Flags

#### Polling Interval

Set clipboard check frequency:

```bash
tracetrim --clipboard-polling-interval 1000ms
```

**Default:** 500ms  
**Impact:** Lower = more responsive, higher CPU. Higher = less responsive, lower CPU.

#### Content Size Limit

Prevent memory exhaustion from huge traces:

```bash
tracetrim --clipboard-max-content-size 5242880  # 5MB
```

**Default:** 1MB (1048576 bytes)

#### Verbose Output

Debug mode with detailed logging:

```bash
tracetrim --verbose
```

Shows:
- Clipboard content changes
- Detection results
- Cleaning operations
- Configuration loaded

#### Quiet Mode

Suppress non-essential output:

```bash
tracetrim --quiet
```

#### Show Timestamps

Include timestamps in output:

```bash
tracetrim --show-timestamp
```

## Configuration File

Create `config.yaml` in working directory:

```yaml
# Clipboard monitoring
clipboard-polling-interval: 500ms
clipboard-max-content-size: 1048576  # 1MB

# Output
verbose: false
quiet: false
show-timestamp: true

# Parser settings
parser-min-stack-lines: 2
parser-min-stack-trace-length: 20

# Script mode
auto-detect-script-mode: true
```

Load custom config:

```bash
tracetrim --config my-config.yaml
```

### Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `clipboard-polling-interval` | 500ms | How often to check clipboard |
| `clipboard-max-content-size` | 1MB | Max clipboard content size |
| `verbose` | false | Enable detailed logging |
| `quiet` | false | Suppress non-essential output |
| `show-timestamp` | true | Show timestamps in output |
| `parser-min-stack-lines` | 2 | Minimum lines for trace detection |
| `parser-min-stack-trace-length` | 20 | Minimum trace length (chars) |
| `auto-detect-script-mode` | true | Auto-detect non-interactive use |

## Examples

### Development Workflow

1. Copy error from browser console
2. TraceTrim detects and cleans trace
3. Paste cleaned trace into issue/chat

```
# Before (cluttered):
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)

# After (clean):
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15) // [x3]
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
```

### CI/CD Integration

Clean stack traces in logs:

```bash
#!/bin/bash
# Extract error logs
grep -A 100 "Error:" app.log | tracetrim --script-mode > errors.txt

# Parse for review
cat errors.txt | grep -E "at |Error"
```

### Debugging Session

Monitor while developing:

```bash
# Terminal 1: Run your app
npm run dev

# Terminal 2: Run TraceTrim
tracetrim --verbose

# Copy errors from browser → auto-cleaned in clipboard
```

### Batch Processing

Clean multiple traces:

```bash
#!/bin/bash
for file in errors/*.log; do
  echo "Cleaning $file..."
  tracetrim --script-mode < "$file" > "cleaned/$(basename $file)"
done
```

## Supported Stack Trace Formats

Detects and cleans:

- **Browser errors:** `Error:`, `TypeError:`, `ReferenceError:`
- **React errors:** React component errors with `react-dom` frames
- **Node.js errors:** Node stack traces with native module frames
- **Custom errors:** Any format with `at functionName (file.js:line:col)` frames

### Example Formats

#### Browser Console

```javascript
Error: Something went wrong
    at myFunction (script.js:10:5)
    at anotherFunction (script.js:15:12)
```

#### React Component Error

```javascript
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
    at MyComponent.render (MyComponent.js:25:10)
```

#### Node.js Error

```
ReferenceError: x is not defined
    at eval (eval at <anonymous> (script.js:1:1))
    at REPL1:1:1
    at Script.runInThisContext (vm.js:122:22)
```

## Troubleshooting

### Traces Not Being Detected

**Symptom:** Copy trace, nothing happens

**Solutions:**
1. Verify clipboard content: `xclip -selection clipboard -o` (Linux)
2. Check detection criteria met:
   - Contains `Error:` or `TypeError:` or similar
   - Has at least 2 stack frames
   - Trace length > 20 characters
3. Enable verbose: `tracetrim --verbose` to see why detection failed

### High CPU Usage

**Symptom:** TraceTrim consuming lots of CPU

**Solutions:**
1. Increase polling interval: `--clipboard-polling-interval 2000ms`
2. Check for rapid clipboard changes (clipboard spam)
3. Reduce max content size limit

### Memory Usage High

**Symptom:** Process consuming >100MB

**Solutions:**
1. Lower content size limit: `--clipboard-max-content-size 524288` (512KB)
2. Restart the process
3. Check for stuck operations in logs

### Script Mode Not Auto-Detecting

**Symptom:** Running in pipe but not in script mode

**Solutions:**
1. Use explicit flag: `--script-mode`
2. Check if TTY is detected: `[ -t 0 ] && echo TTY || echo no-TTY`
3. Set `auto-detect-script-mode: true` in config

## Best Practices

- **Keep running:** Leave tracetrim running in background during development
- **Polling interval:** Higher interval for low-power devices
- **Size limits:** Adjust based on your stack trace sizes
- **Verbose mode:** Enable when debugging detection issues
- **Config file:** Create `config.yaml` for persistent settings

## Integration with Tools

### IDE Plugins

Some IDEs may support clipboard integration:
- VSCode: Copy stack trace → TraceTrim → paste cleaned
- JetBrains: Clipboard monitor integration (if available)

### Error Tracking Services

Clean traces before sending to Sentry/Rollbar:

```bash
# Extract trace, clean it, send
tracetrim --script-mode < raw_trace.txt | curl -X POST \
  https://sentry.io/api/... \
  -d @-
```

## Performance

Typical performance:
- **CPU:** <2% at 500ms polling
- **Memory:** <50MB at rest
- **Parsing:** <50ms for 100-frame traces
- **Startup:** <100ms

## Support

For issues:
1. Check verbose output: `--verbose`
2. File issue with: tracetrim version, config, example trace, error message
3. Enable full logging for debugging

---

**Last updated:** 2026-04-28
