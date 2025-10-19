# TraceTrim

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rethunk-tech/tracetrim)](https://goreportcard.com/report/github.com/rethunk-tech/tracetrim)

[![Release](https://img.shields.io/github/v/release/rethunk-tech/tracetrim.svg)](https://github.com/rethunk-tech/tracetrim/releases)
[![Downloads](https://img.shields.io/github/downloads/rethunk-tech/tracetrim/total.svg)](https://github.com/rethunk-tech/tracetrim/releases)
[![Platforms](https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-blue.svg)](https://github.com/rethunk-tech/tracetrim)
[![Issues](https://img.shields.io/github/issues/rethunk-tech/tracetrim.svg)](https://github.com/rethunk-tech/tracetrim/issues)
[![Stars](https://img.shields.io/github/stars/rethunk-tech/tracetrim.svg)](https://github.com/rethunk-tech/tracetrim)

A simple, cross-platform CLI application that monitors your clipboard for JavaScript console or React stack traces and automatically trims them by removing repetitive blocks.

## Problem Solved

React stack traces often contain repetitive blocks of text that make them hard to read. For example:

**Before ( cluttered with repetitive frames):**

```console
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
```

**After (clean and readable):**

```console
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15) // [x4]
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
```

## Features

- 🚀 **Automatic Detection**: Continuously monitors clipboard for stack traces
- 🎯 **Smart Cleaning**: Removes only repetitive blocks, preserves all formatting
- ⚡ **Real-time**: Updates clipboard instantly when stack traces are detected
- 🔧 **Script Mode**: Can be used in shell scripts and automation pipelines
- 🌍 **Cross-platform**: Works on Windows, macOS, and Linux using golang.design/x/clipboard
- 🧪 **Well-tested**: Comprehensive test coverage for reliable operation
- 📦 **Zero-config**: Just run it - no configuration needed

## Installation

### Option 1: Download Pre-built Binaries

Check the [Releases](https://github.com/rethunk-tech/tracetrim/releases) page for pre-built binaries for your platform.

### Option 2: Build from Source

**Prerequisites:**

- Go 1.21 or later

```bash
# Clone the repository
git clone https://github.com/rethunk-tech/tracetrim.git
cd tracetrim

# Build for your platform
go build -o tracetrim ./cmd/

# Optional: Install to system PATH
# Linux/macOS:
# sudo mv tracetrim /usr/local/bin/
# Windows:
# move tracetrim %PATH%
```

### Option 3: Cross-platform Builds

The application uses a cross-platform clipboard library, so you can build for multiple platforms:

```bash
# Build for multiple platforms
GOOS=windows GOARCH=amd64 go build -o tracetrim-windows.exe ./cmd/
GOOS=darwin GOARCH=amd64 go build -o tracetrim-macos ./cmd/
GOOS=linux GOARCH=amd64 go build -o tracetrim-linux ./cmd/

# Build for ARM architectures
GOOS=darwin GOARCH=arm64 go build -o tracetrim-macos-arm64 ./cmd/
GOOS=linux GOARCH=arm64 go build -o tracetrim-linux-arm64 ./cmd/
```

## Usage

### Basic Usage

Simply run the application:

```bash
./tracetrim
```

The application will:

1. Start monitoring your clipboard
2. Display a message indicating it's running
3. Automatically detect and clean stack traces when you copy them
4. Show a brief notification when cleaning occurs

**Example output:**

```console
TraceTrim
Monitoring clipboard for JavaScript/React stack traces...
Press Ctrl+C to exit

[14:23:45] Detected stack trace, cleaning...
✓ Stack trace cleaned and clipboard updated
  Removed 3 repetitive lines
```

### Stopping the Application

Press `Ctrl+C` to gracefully stop the clipboard monitoring.

### Script Mode

TraceTrim can be used in scripts in two ways:

#### Automatic Detection (Recommended)

TraceTrim automatically detects when it's being used in a non-interactive environment (scripts, pipes, redirection) and switches to script mode automatically:

```bash
# Clean a stack trace from a file (auto-detected)
cat stack_trace.txt | ./tracetrim > cleaned_stack_trace.txt

# Use in a pipeline (auto-detected)
echo "Error: Something went wrong
    at myFunction (script.js:10:5)
    at myFunction (script.js:10:5)" | ./tracetrim

# Process multiple stack traces (auto-detected)
find . -name "*.log" -exec grep -l "Error:" {} \; | xargs cat | ./tracetrim
```

#### Manual Script Mode

You can also explicitly enable script mode with the `--script-mode` flag:

```bash
# Explicitly enable script mode
cat stack_trace.txt | ./tracetrim --script-mode > cleaned_stack_trace.txt
```

**Script Mode Features:**

- **Automatic Detection**: Detects pipes, redirection, and CI environments
- Reads stack traces from STDIN
- Outputs cleaned content to STDOUT
- No headers, footers, or status messages (clean output for scripts)
- Exits immediately after processing
- Compatible with shell pipelines and automation scripts

## Configuration

TraceTrim supports both configuration files and command-line flags. By default, it looks for a `config.yaml` file in the current directory.

### Command Line Flags

- `--verbose`: Enable detailed logging output
- `--quiet`: Suppress non-essential output
- `--script-mode`: Enable script mode (overrides auto-detection)
- `--auto-detect-script-mode`: Auto-detect script mode (default: true)
- `--clipboard-polling-interval`: Set clipboard polling interval (default: 500ms)
- `--clipboard-max-content-size`: Maximum clipboard content size in bytes (default: 1MB)
- `--parser-min-stack-lines`: Minimum stack lines for detection (default: 2)
- `--parser-min-stack-trace-length`: Minimum stack trace length (default: 20)
- `--show-timestamp`: Show timestamps in output (default: true)
- `--config`: Specify configuration file path (default: config.yaml)

### Configuration File

Example `config.yaml`:

```yaml
# Clipboard monitoring settings
clipboard-polling-interval: 500ms
clipboard-max-content-size: 1048576

# Output settings
verbose: false
quiet: false
show-timestamp: true

# Parser settings
parser-min-stack-lines: 2
parser-min-stack-trace-length: 20

# Script mode settings
auto-detect-script-mode: true
```

## How It Works

### Clipboard Integration

The application uses the [golang.design/x/clipboard](https://github.com/golang.design/x/clipboard) library for cross-platform clipboard access, providing seamless support for Windows, macOS, and Linux without platform-specific code.

### Detection

The application uses pattern matching to identify JavaScript and React stack traces in clipboard content. It looks for:

- JavaScript error patterns (`Error:`, `TypeError:`, `ReferenceError:`)
- React-specific patterns (`react-dom.development.js`, `ReactErrorUtils.invokeGuardedCallback`)
- Stack frame patterns (`at functionName (file.js:line:column)`)

### Cleaning Process

1. **Parse**: Analyze the clipboard content to identify stack trace patterns
2. **Identify Duplicates**: Find repetitive stack frames (same function, file, and line)
3. **Remove**: Eliminate duplicate frames while preserving the first occurrence
4. **Preserve Formatting**: Maintain all original indentation, spacing, and non-stack-trace content
5. **Update**: Replace clipboard content with the cleaned version
6. **Notify**: Show a brief message indicating what was cleaned

### What Gets Preserved

- ✅ Original error messages
- ✅ All unique stack frames
- ✅ Indentation and formatting
- ✅ Non-stack-trace content
- ✅ File paths and line numbers

### What Gets Removed

- ❌ Duplicate stack frames (same function + file + line)
- ❌ Nothing else - all formatting and content is preserved

## Supported Stack Trace Formats

The application handles various JavaScript and React stack trace formats:

```javascript
// Browser console errors
Error: Something went wrong
    at myFunction (script.js:10:5)
    at anotherFunction (script.js:15:12)

// React component errors
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
    at MyComponent.render (MyComponent.js:25:10)

// Node.js errors
ReferenceError: x is not defined
    at eval (eval at <anonymous> (script.js:1:1))
```

## Security Considerations

This application handles clipboard content and should be used with appropriate security considerations:

### Security Features

- **Input Validation**: All clipboard content is validated before processing
  - UTF-8 validation ensures no binary data corruption
  - Size limits prevent memory exhaustion attacks
  - Content sanitization removes potentially dangerous patterns
- **Memory Safety**: Proper memory management prevents leaks and corruption
- **Platform Compatibility**: Cross-platform library handles platform differences automatically
- **No Network Access**: Application operates entirely locally with no external connections
- **Minimal Permissions**: Only requires clipboard access permissions

### Security Best Practices

- **Permission Management**: Standard application permissions for clipboard access
- **Content Safety**:
  - Application only processes text content, never binary data
  - Stack trace patterns are strictly validated before processing
  - No execution of clipboard content or pattern injection
- **Resource Limits**:
  - Configurable content size limits (default: 1MB)
  - Processing timeout limits prevent hanging
  - Memory usage is bounded and monitored

### Potential Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Memory Exhaustion** | High | Content size limits, input validation, bounded processing |
| **Pattern Injection** | Medium | Strict regex validation, no dynamic pattern execution |
| **Clipboard Pollution** | Low | Content validation, safe fallback on errors |
| **Resource Leaks** | Medium | Proper cleanup, defer statements, error handling |
| **Race Conditions** | Low | Synchronization primitives, atomic operations |

### Threat Model

- **Attack Vectors Considered**:
  - Malicious clipboard content designed to crash the application
  - Extremely large content causing memory exhaustion
  - Invalid UTF-8 sequences causing parsing errors
  - Rapid clipboard changes causing race conditions
- **Attack Vectors Not Applicable**:
  - Network-based attacks (no network access)
  - File system attacks (no file I/O beyond config)
  - Privilege escalation (runs with user permissions only)
  - Code injection (no dynamic code execution)

### Security Updates

The application follows security best practices:

- Regular dependency updates through `go mod tidy`
- Static analysis with `go vet` to catch potential issues
- Comprehensive test coverage including edge cases
- Cross-platform security validations

If you discover any security vulnerabilities, please report them responsibly through the project's issue tracker.

## Supported File Types

The application supports cleaning stack traces from all modern JavaScript and TypeScript file formats:

- **JavaScript (.js)** - Standard JavaScript files
- **TypeScript (.ts)** - TypeScript source files
- **React JSX (.jsx)** - React JavaScript XML components
- **React TSX (.tsx)** - React TypeScript XML components
- **Modern modules (.mjs)** - ES modules and modern JavaScript

The parser automatically detects and handles stack traces from any of these file types, regardless of build tools, bundlers, or development environments being used.

## Troubleshooting Guide

This section provides comprehensive solutions for common issues you might encounter when using the TraceTrim.

### Quick Diagnostic Commands

Before diving into specific issues, try these diagnostic commands:

```bash
# Check if clipboard tools are available (Linux)
which xclip xsel

# Check X11 display (Linux)
echo $DISPLAY

# Check Go version compatibility
go version

# Test basic clipboard functionality
echo "test content" | xclip -selection clipboard -i  # Linux
echo "test content" | pbcopy  # macOS
```

### Common Issues and Solutions

#### 1. Application Won't Start

**Symptoms**: Application fails to start or exits immediately
**Possible Causes**: Missing dependencies, permission issues, corrupted binary

**Solutions**:

- **Check Go installation**: Ensure Go 1.21+ is installed (`go version`)
- **Verify binary**: Ensure the binary is not corrupted (`ls -la tracetrim`)
- **Check permissions**: Ensure the binary is executable (`chmod +x tracetrim`)
- **Platform compatibility**: Ensure you're running on a supported platform (Windows, macOS, Linux)

#### 2. Clipboard Access Denied

**Symptoms**: "Failed to initialize clipboard monitor" errors
**Possible Causes**: Missing permissions, incompatible environment, system restrictions

**Solutions**:

- Ensure the application has permission to access the clipboard
- Try restarting the application
- Check if another application has exclusive clipboard access
- On Linux, ensure a clipboard manager is running

#### 3. Stack Traces Not Being Cleaned

**Symptoms**: Stack traces appear in clipboard but are not processed
**Possible Causes**: Detection issues, content format problems, configuration issues

**Solutions**:

- **Verify content format**: Ensure the clipboard contains actual JavaScript/React stack traces
- **Check content size**: Very large stack traces (>1MB) are skipped by default
- **Enable verbose mode**: Use `--verbose` flag to see detailed processing information
- **Adjust polling interval**: Use `--clipboard-polling-interval` flag (default: 500ms)
- **Check for false negatives**: Some minified or non-standard stack trace formats may not be detected

#### 4. High CPU Usage

**Symptoms**: Application consumes excessive CPU resources
**Possible Causes**: Very frequent clipboard polling, large content processing

**Solutions**:

- **Adjust polling interval**: Use `--clipboard-polling-interval` flag (default: 500ms)
- **Check for clipboard spam**: Rapid clipboard changes can cause high CPU usage
- **Monitor resource usage**: Use system tools to identify the cause

#### 5. Memory Issues

**Symptoms**: Application crashes or uses excessive memory
**Possible Causes**: Large clipboard content, memory leaks, system resource constraints

**Solutions**:

- **Check content size limits**: Large content (>1MB) is skipped by default
- **Check system resources**: Ensure adequate free memory (at least 50MB recommended)
- **Monitor for memory leaks**: Use system tools to check application memory usage

#### 7. Configuration Issues

**Symptoms**: Settings not taking effect, configuration file errors
**Possible Causes**: Invalid configuration file, permission issues, syntax errors

**Solutions**:

- **Validate config file**: Check `config.yaml` for syntax errors
- **Check file permissions**: Ensure config file is readable
- **Reset configuration**: Delete `config.yaml` to use defaults
- **Command line override**: Use command line flags to override config file settings

### Getting Help

If you encounter issues not covered in this guide:

1. **Check the logs**: Run with `--verbose` flag for detailed information
2. **Review recent changes**: Check if issues started after system/application updates
3. **Test with simple content**: Try copying simple text to verify basic clipboard functionality
4. **Check system requirements**: Verify your system meets minimum requirements
5. **Report issues**: Use the project's issue tracker with detailed information including:
   - Platform and version
   - Go version
   - Error messages
   - Steps to reproduce
   - Relevant system information

### Advanced Debugging

For advanced users experiencing persistent issues:

```bash
# Enable verbose logging
./tracetrim --verbose

# Test clipboard access directly
./tracetrim --clipboard-polling-interval 2000ms --verbose

# Check system clipboard status
# Use platform-specific clipboard tools (pbpaste on macOS, xclip/xsel on Linux, Get-Clipboard on Windows)

# Monitor application with system tools
# Use platform-specific monitoring tools (Activity Monitor on macOS, Task Manager on Windows, top/htop on Linux)
```

### False Positives

The application is designed to be conservative - it will only clean content that clearly matches stack trace patterns. If you encounter false positives, please report them as issues.

### Performance

The application uses minimal system resources:

- Polls clipboard every 500ms
- Uses efficient string processing
- Minimal memory footprint

## Development

### Running Tests

```bash
go test ./...
```

### Code Structure

```tree
├── cmd/                    # CLI application entry point
│   └── main.go
├── clipboard/              # Clipboard monitoring module
│   ├── monitor.go          # Cross-platform clipboard interface using golang.design/x/clipboard
│   └── monitor_test.go     # Clipboard monitoring tests
├── parser/                 # Stack trace parsing and cleaning
│   ├── parser.go           # Core parsing logic
│   └── parser_test.go      # Comprehensive tests
├── internal/
│   ├── config/             # Configuration management
│   │   ├── config.go        # Configuration loading and validation
│   │   └── config_test.go  # Configuration tests
│   └── models/             # Shared data structures
│       └── types.go
└── go.mod                  # Go module definition
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

### Areas for Enhancement

- Additional stack trace formats
- Configuration options
- GUI interface
- Integration with editors/IDEs

## License

This project is open source. See LICENSE file for details.

## Support

If you find this tool helpful, please consider:

- ⭐ Starring the repository
- 🐛 Reporting bugs or issues
- 💡 Suggesting improvements
- 🚀 Contributing code

---

[LICENSE](LICENSE) | Copyright (c) 2025 Rethunk.Tech, LLC.
