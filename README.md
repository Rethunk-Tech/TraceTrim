# Clipboard Stack Trace Cleaner

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/rethunk-tech/no-reaction)](https://goreportcard.com/report/github.com/rethunk-tech/no-reaction)

[![Release](https://img.shields.io/github/v/release/rethunk-tech/no-reaction.svg)](https://github.com/rethunk-tech/no-reaction/releases)
[![Downloads](https://img.shields.io/github/downloads/rethunk-tech/no-reaction/total.svg)](https://github.com/rethunk-tech/no-reaction/releases)
[![Platforms](https://img.shields.io/badge/platforms-Windows%20%7C%20macOS%20%7C%20Linux-blue.svg)]()
[![Issues](https://img.shields.io/github/issues/rethunk-tech/no-reaction.svg)](https://github.com/rethunk-tech/no-reaction/issues)
[![Stars](https://img.shields.io/github/stars/rethunk-tech/no-reaction.svg)](https://github.com/rethunk-tech/no-reaction)

A simple, cross-platform CLI application that monitors your clipboard for JavaScript console or React stack traces and automatically cleans them by removing repetitive blocks.

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
// Removed 3 repetitive stack frame(s)
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
```

## Features

- 🚀 **Automatic Detection**: Continuously monitors clipboard for stack traces
- 🎯 **Smart Cleaning**: Removes only repetitive blocks, preserves all formatting
- ⚡ **Real-time**: Updates clipboard instantly when stack traces are detected
- 🌍 **Cross-platform**: Works on Windows, macOS, and Linux
- 🧪 **Well-tested**: Comprehensive test coverage for reliable operation
- 📦 **Zero-config**: Just run it - no configuration needed

## Installation

### Option 1: Download Pre-built Binaries

Check the [Releases](https://github.com/rethunk-tech/no-reaction/releases) page for pre-built binaries for your platform.

### Option 2: Build from Source

**Prerequisites:**

- Go 1.19 or later

```bash
# Clone the repository
git clone https://github.com/rethunk-tech/no-reaction.git
cd no-reaction

# Build for your platform
go build -o clipboard-cleaner ./cmd/

# Optional: Install to system PATH
# Linux/macOS:
# sudo mv clipboard-cleaner /usr/local/bin/
# Windows:
# move clipboard-cleaner %PATH%
```

### Option 3: Cross-platform Builds

The application supports all major platforms with platform-specific optimizations:

```bash
# Build for multiple platforms (build tags ensure correct implementation)
GOOS=windows GOARCH=amd64 go build -o clipboard-cleaner-windows.exe ./cmd/
GOOS=darwin GOARCH=amd64 go build -o clipboard-cleaner-macos ./cmd/
GOOS=linux GOARCH=amd64 go build -o clipboard-cleaner-linux ./cmd/

# Build for ARM architectures
GOOS=darwin GOARCH=arm64 go build -o clipboard-cleaner-macos-arm64 ./cmd/
GOOS=linux GOARCH=arm64 go build -o clipboard-cleaner-linux-arm64 ./cmd/
```

**Note**: Each platform uses its optimal clipboard access method:

- **Windows**: Native Windows API (`user32.dll`, `kernel32.dll`)
- **macOS**: Native Cocoa NSPasteboard (Objective-C bridge via cgo)
- **Linux**: xclip/xsel utilities (automatically detected and fallback)

## Usage

### Basic Usage

Simply run the application:

```bash
./clipboard-cleaner
```

The application will:

1. Start monitoring your clipboard
2. Display a message indicating it's running
3. Automatically detect and clean stack traces when you copy them
4. Show a brief notification when cleaning occurs

**Example output:**

```console
Clipboard Stack Trace Cleaner
Monitoring clipboard for JavaScript/React stack traces...
Press Ctrl+C to exit

[14:23:45] Detected stack trace, cleaning...
✓ Stack trace cleaned and clipboard updated
  Removed 3 repetitive lines
```

### Stopping the Application

Press `Ctrl+C` to gracefully stop the clipboard monitoring.

## How It Works

### Platform Architecture

The application uses a clean platform abstraction architecture:

- **Interface-based Design**: `Platform` interface defines clipboard operations (`GetContent()`, `SetContent()`, `GetName()`)
- **Build Tag System**: Platform-specific implementations are selected at compile time using Go build tags
- **Runtime Detection**: Each platform implementation handles its own initialization and error handling
- **Cross-platform Compatibility**: Single binary works across all platforms with optimal native integration

### Clipboard Integration

Each platform uses its most efficient clipboard access method:

- **Windows**: Direct Windows API calls for maximum performance
- **macOS**: Native Cocoa NSPasteboard integration via Objective-C bridge
- **Linux**: External clipboard utilities (xclip/xsel) with automatic fallback

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
- **Platform Isolation**: Each platform implementation is isolated and validated
- **No Network Access**: Application operates entirely locally with no external connections
- **Minimal Permissions**: Only requires clipboard access permissions

### Security Best Practices

- **Permission Management**:
  - Windows: Standard application permissions (no administrator required)
  - macOS: Requires Accessibility permissions (see troubleshooting section)
  - Linux: Requires clipboard utility access (xclip/xsel)
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
- Platform-specific security validations

If you discover any security vulnerabilities, please report them responsibly through the project's issue tracker.

## Supported File Types

The application supports cleaning stack traces from all modern JavaScript and TypeScript file formats:

- **JavaScript (.js)** - Standard JavaScript files
- **TypeScript (.ts)** - TypeScript source files
- **React JSX (.jsx)** - React JavaScript XML components
- **React TSX (.tsx)** - React TypeScript XML components
- **Modern modules (.mjs)** - ES modules and modern JavaScript

The parser automatically detects and handles stack traces from any of these file types, regardless of build tools, bundlers, or development environments being used.

## Platform-Specific Implementation Details

### Windows

**Implementation**: Native Windows API integration using `user32.dll` and `kernel32.dll`

- **Clipboard Access**: Direct Windows clipboard API calls for optimal performance
- **Requirements**: Windows Vista or later (Windows 7+ recommended)
- **Dependencies**: None - uses only standard Windows libraries
- **Permissions**: Standard application permissions, no special setup required
- **Build Tags**: `//go:build windows` for platform-specific compilation

**Technical Details**:

- Uses `OpenClipboard`, `GetClipboardData`, `SetClipboardData` Windows APIs
- Handles UTF-16 text conversion for Windows clipboard format
- Proper memory management with `GlobalAlloc`/`GlobalFree` and `GlobalLock`/`GlobalUnlock`

### macOS

**Implementation**: Cocoa NSPasteboard integration using Objective-C bridge via cgo

- **Clipboard Access**: Native macOS NSPasteboard API for seamless integration
- **Requirements**: macOS 10.6 or later (macOS 10.15+ recommended)
- **Dependencies**: None - uses only system Cocoa framework
- **Permissions**: May require Accessibility permissions on first run
- **Build Tags**: `//go:build darwin` for platform-specific compilation

**Technical Details**:

- Cgo interface to Objective-C Cocoa NSPasteboard APIs
- Automatic memory management with `@autoreleasepool`
- Handles NSString conversion between Go and Objective-C
- Uses `NSPasteboardTypeString` for text content

**Setup**:

```bash
# Grant permissions if prompted
# System Preferences > Security & Privacy > Accessibility
```

### Linux

**Implementation**: External clipboard utilities with fallback support

- **Primary Tool**: xclip (supports both X11 and Wayland via XWayland)
- **Fallback Tool**: xsel (alternative clipboard utility)
- **Requirements**: X11 environment or Wayland with XWayland
- **Dependencies**: Install either xclip or xsel (automatically detected)

**Technical Details**:

- Uses `exec.Command` to interface with external clipboard tools
- Automatic tool detection and fallback mechanism
- Supports both reading (`xclip -o`/`xsel -ob`) and writing (`xclip -i`/`xsel -ib`)
- Handles text content via stdin/stdout pipes

**Installation**:

```bash
# Ubuntu/Debian
sudo apt-get install xclip

# Alternative for Ubuntu/Debian
sudo apt-get install xsel

# CentOS/RHEL/Fedora
sudo yum install xsel

# Arch Linux
sudo pacman -S xclip
```

**Environment Support**:

- ✅ X11 desktop environments (GNOME, KDE, Xfce, etc.)
- ✅ Wayland with XWayland compatibility layer
- ✅ Remote X11 sessions via SSH with X forwarding

## Troubleshooting Guide

This section provides comprehensive solutions for common issues you might encounter when using the Clipboard Stack Trace Cleaner.

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
- **Verify binary**: Ensure the binary is not corrupted (`ls -la clipboard-cleaner`)
- **Check permissions**: Ensure the binary is executable (`chmod +x clipboard-cleaner`)
- **Platform compatibility**: Ensure you're running on a supported platform (Windows 7+, macOS 10.6+, Linux with X11)

#### 2. Clipboard Access Denied

**Symptoms**: "Failed to initialize clipboard monitor" errors
**Possible Causes**: Missing permissions, incompatible environment, system restrictions

**Platform-Specific Solutions**:

**Windows**:

- Run as administrator if clipboard access fails
- Check Windows version (Vista+ required)
- Verify no other application has exclusive clipboard access
- Try running from Command Prompt instead of PowerShell

**macOS**:

- Grant Accessibility permissions (see detailed steps in Platform-Specific section)
- Check System Preferences > Security & Privacy > Privacy > Accessibility
- Ensure terminal application (Terminal, iTerm, etc.) has permissions
- Try restarting the application after granting permissions

**Linux**:

- Install clipboard utilities: `sudo apt-get install xclip` (Ubuntu/Debian) or `sudo yum install xsel` (CentOS/RHEL)
- Verify X11 is running: `echo $DISPLAY` should show a display number
- Check if in SSH session without X forwarding: use `ssh -X` for X11 forwarding

#### 3. Stack Traces Not Being Cleaned

**Symptoms**: Stack traces appear in clipboard but are not processed
**Possible Causes**: Detection issues, content format problems, configuration issues

**Solutions**:

- **Verify content format**: Ensure the clipboard contains actual JavaScript/React stack traces
- **Check content size**: Very large stack traces (>1MB) are skipped by default
- **Adjust detection sensitivity**: Use `--parser-min-stack-lines` flag to lower detection threshold
- **Enable verbose mode**: Use `--verbose` flag to see detailed processing information
- **Check for false negatives**: Some minified or non-standard stack trace formats may not be detected

#### 4. High CPU Usage

**Symptoms**: Application consumes excessive CPU resources
**Possible Causes**: Very frequent clipboard polling, large content processing

**Solutions**:

- **Increase polling interval**: Use `--clipboard-polling-interval 1000ms` (default: 500ms)
- **Check for clipboard spam**: Rapid clipboard changes can cause high CPU usage
- **Monitor resource usage**: Use system tools to identify the cause

#### 5. Memory Issues

**Symptoms**: Application crashes or uses excessive memory
**Possible Causes**: Large clipboard content, memory leaks, system resource constraints

**Solutions**:

- **Configure content limits**: Use `--clipboard-max-content-size 1048576` (1MB limit)
- **Check system resources**: Ensure adequate free memory (at least 50MB recommended)
- **Monitor for memory leaks**: Use system tools to check application memory usage

#### 6. Cross-Platform Issues

**Symptoms**: Application works on one platform but not others
**Possible Causes**: Platform-specific dependencies, compilation issues, environment differences

**Solutions**:

- **Verify platform support**: Check you're using a supported OS/architecture combination
- **Check compilation**: Ensure binary was compiled for the correct platform
- **Environment variables**: Check for platform-specific environment requirements
- **Dependencies**: Verify all platform-specific dependencies are installed

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
./clipboard-cleaner --verbose

# Test clipboard access directly
./clipboard-cleaner --clipboard-polling-interval 2000ms

# Check system clipboard status
# Windows: Use ClipBook Viewer or PowerShell Get-Clipboard
# macOS: Use pbpaste
# Linux: Use xclip -o or xsel -ob

# Monitor application with system tools
# Windows: Task Manager, Resource Monitor
# macOS: Activity Monitor, Console.app
# Linux: top, htop, journalctl
```

## Platform-Specific Implementation Details

### Clipboard Access Issues

**Windows:**

- **Permission Issues**: Ensure the application has permission to access clipboard (standard Windows permissions)
- **Administrator Rights**: Try running as administrator if clipboard access fails
- **Windows Version**: Requires Windows Vista or later; Windows 7+ recommended
- **Error Messages**: Check for Windows API errors in console output

**macOS:**

- **Accessibility Permissions**: Grant permissions in System Preferences > Security & Privacy > Accessibility
  - Go to System Preferences > Security & Privacy > Privacy tab
  - Select "Accessibility" from the left sidebar
  - Click the lock icon to make changes (you may need to enter your password)
  - Find and check the box next to your terminal application (Terminal, iTerm, etc.) or the clipboard-cleaner executable
  - If the application doesn't appear in the list, try adding it manually by clicking the "+" button
- **First Run Prompt**: The application may prompt for permissions on first launch
- **System Integrity Protection**: Ensure SIP doesn't interfere with clipboard access
  - SIP should not normally interfere with clipboard access, but if you encounter issues, check that the application is properly signed
- **Error Messages**: Look for cgo/Objective-C bridge errors in console output
- **Troubleshooting Steps**:
  1. Ensure the application has Accessibility permissions as described above
  2. Try running the application from a terminal that has proper permissions
  3. If using a bundled application, ensure it's properly signed with codesign
  4. Check Console.app for any related error messages

**Linux:**

- **Missing Dependencies**: Install xclip or xsel (see installation section above)
- **X11 Environment**: Ensure X11 is running (`echo $DISPLAY` should show display)
- **Wayland Compatibility**: Use XWayland for Wayland environments
- **SSH Sessions**: Enable X11 forwarding (`ssh -X`) for remote sessions
- **Permission Issues**: May need `xhost +` for local clipboard access
- **Tool Detection**: Application automatically detects available tools and falls back

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
│   ├── monitor.go          # Cross-platform clipboard interface
│   ├── monitor_windows.go  # Windows-specific implementation
│   ├── monitor_darwin.go   # macOS-specific implementation (Cocoa NSPasteboard)
│   └── monitor_linux.go    # Linux-specific implementation (xclip/xsel)
├── parser/                 # Stack trace parsing and cleaning
│   ├── parser.go           # Core parsing logic
│   └── parser_test.go      # Comprehensive tests
└── internal/models/        # Shared data structures
    └── types.go
```

### Adding New Platforms

To add support for a new platform:

1. Create `monitor_[platform].go` in the `clipboard/` directory
2. Add appropriate build tags (`//go:build [platform]` and `// +build [platform]`)
3. Implement the `Platform` interface with `GetContent()`, `SetContent()`, and `GetName()` methods
4. Add a `getPlatform()` function that returns your platform implementation
5. Update this documentation to include platform-specific details

**Platform Interface:**

```go
type Platform interface {
    GetContent() (string, error)
    SetContent(content string) error
    GetName() string
}
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
