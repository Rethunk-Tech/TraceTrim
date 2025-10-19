# Clipboard Stack Trace Cleaner

A simple Golang CLI application that monitors the clipboard for JavaScript console or React stack traces and automatically cleans them by removing repetitive blocks.

## Purpose

React stack traces often contain repetitive blocks of text that make them hard to read. This application will:
1. Monitor clipboard changes continuously
2. Detect JavaScript/React stack traces in clipboard content
3. Remove repetitive blocks while preserving the essential error information
4. Replace the clipboard content with the cleaned version

## Implementation Steps

### Phase 1: Project Setup
1. **Initialize Git Repository** - Set up version control
2. **Create Project Structure** - Organize code into logical modules
3. **Set up Go Module** - Initialize go.mod

### Phase 2: Core Modules

#### 1. Clipboard Monitor (`clipboard/monitor.go`)
**Purpose**: Continuously watch for clipboard changes
- Use platform-specific clipboard APIs (Windows: win32 API, macOS: Cocoa NSPasteboard, Linux: xclip/xsel)
- Poll clipboard at reasonable intervals (e.g., 500ms)
- Trigger processing when content changes
- Handle platform differences gracefully

**Key Functions**:
- `StartMonitoring()` - Begin clipboard monitoring
- `StopMonitoring()` - Stop monitoring
- `GetCurrentContent()` - Get current clipboard content
- `SetContent(string)` - Set clipboard content

#### 2. Stack Trace Parser (`parser/parser.go`)
**Purpose**: Detect and clean stack traces
- Identify JavaScript/React stack traces in clipboard content
- Remove repetitive blocks (common in React error handling)
- Preserve essential error information
- Handle various stack trace formats

**Key Functions**:
- `IsStackTrace(string) bool` - Detect if content is a stack trace
- `CleanStackTrace(string) string` - Remove repetitive blocks
- `ExtractErrorInfo(string) ErrorInfo` - Extract key error details

#### 3. Main CLI Application (`cmd/main.go`)
**Purpose**: Tie everything together in a simple CLI interface
- Initialize clipboard monitoring
- Handle graceful shutdown (Ctrl+C)
- Provide minimal user feedback
- Keep it simple - no complex argument parsing

### Phase 3: Testing
1. **Unit Tests** for each module
2. **Integration Tests** for end-to-end functionality
3. **Platform-specific Tests** for clipboard operations

### Phase 4: Documentation
1. **README.md** with usage instructions
2. **Build Instructions** for different platforms
3. **Examples** of before/after stack trace cleaning

## Technical Considerations

### Platform Support
- **Windows**: Use Windows API for clipboard access
- **macOS**: Use Cocoa NSPasteboard
- **Linux**: Use xclip or xsel
- **Cross-platform**: Abstract platform differences

### Stack Trace Patterns
React stack traces typically repeat blocks like:
```
at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
```
Should be reduced to single occurrence while preserving the error context.

### Error Handling
- Graceful handling of clipboard access failures
- Platform-specific error messages
- Don't crash on malformed input

### Performance
- Minimal resource usage
- Efficient string processing
- Reasonable polling interval

## File Structure
```
clipboard-cleaner/
├── cmd/
│   └── main.go              # CLI entry point
├── clipboard/
│   ├── monitor.go           # Clipboard monitoring
│   └── monitor_windows.go   # Windows-specific implementation
├── parser/
│   └── parser.go            # Stack trace detection and cleaning
├── internal/
│   └── models/
│       └── types.go         # Shared types and structs
├── go.mod
├── README.md
└── PROJECT_PLAN.md
```

## Success Criteria
- ✅ Monitors clipboard continuously
- ✅ Detects JavaScript/React stack traces
- ✅ Removes repetitive blocks
- ✅ Replaces clipboard content automatically
- ✅ Works across platforms
- ✅ Has comprehensive tests
- ✅ Simple CLI interface
- ✅ Minimal resource usage

## Future Enhancements (Out of Scope)
- Configuration file for custom patterns
- GUI interface
- Multiple clipboard formats support
- Stack trace syntax highlighting
- Integration with editors/IDEs
