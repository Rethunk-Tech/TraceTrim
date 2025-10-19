# Clipboard Stack Trace Cleaner

A simple, cross-platform CLI application that monitors your clipboard for JavaScript console or React stack traces and automatically cleans them by removing repetitive blocks.

## Problem Solved

React stack traces often contain repetitive blocks of text that make them hard to read. For example:

**Before ( cluttered with repetitive frames):**
```
Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
```

**After (clean and readable):**
```
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

```bash
# Build for multiple platforms
GOOS=windows GOARCH=amd64 go build -o clipboard-cleaner-windows.exe ./cmd/
GOOS=darwin GOARCH=amd64 go build -o clipboard-cleaner-macos ./cmd/
GOOS=linux GOARCH=amd64 go build -o clipboard-cleaner-linux ./cmd/
```

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
```
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

## Platform-Specific Notes

### Windows
- Uses Windows API for clipboard access
- Requires Windows Vista or later
- No additional dependencies

### macOS
- Uses Cocoa NSPasteboard API
- Requires macOS 10.6 or later
- No additional dependencies

### Linux
- Uses xclip or xsel (automatically detected)
- Install xclip: `sudo apt-get install xclip` (Ubuntu/Debian)
- Install xsel: `sudo yum install xsel` (CentOS/RHEL)
- Fallback: Direct X11 clipboard access

## Troubleshooting

### Clipboard Access Issues

**Windows:**
- Ensure the application has permission to access clipboard
- Try running as administrator if issues persist

**macOS:**
- Grant Accessibility permissions in System Preferences > Security & Privacy
- The application may prompt for permissions on first run

**Linux:**
- Install xclip or xsel as mentioned above
- Ensure X11 is running and accessible

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

```
├── cmd/                    # CLI application entry point
│   └── main.go
├── clipboard/              # Clipboard monitoring module
│   ├── monitor.go          # Cross-platform clipboard interface
│   └── monitor_windows.go  # Windows-specific implementation
├── parser/                 # Stack trace parsing and cleaning
│   ├── parser.go           # Core parsing logic
│   └── parser_test.go      # Comprehensive tests
└── internal/models/        # Shared data structures
    └── types.go
```

### Adding New Platforms

To add support for a new platform:

1. Create `monitor_[platform].go` in the `clipboard/` directory
2. Implement the `Platform` interface
3. Update `getPlatform()` in `monitor.go` to detect and return your implementation

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
