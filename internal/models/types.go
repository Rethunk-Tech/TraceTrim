package models

import "time"

// ErrorInfo contains the essential information from a stack trace
type ErrorInfo struct {
	Message   string   // The main error message
	Stack     []string // Individual stack frames
	Source    string   // Original source file/line if available
	Component string   // React component name if applicable
}

// ClipboardContent represents clipboard data with metadata
type ClipboardContent struct {
	Content   string    // The actual clipboard content
	Timestamp time.Time // When this content was captured
	Format    string    // Content format (text/plain, etc.)
}

// StackFrame represents a single frame in a stack trace
type StackFrame struct {
	Function string // Function name
	File     string // File path
	Line     int    // Line number
	Column   int    // Column number (if available)
}

// CleanResult contains the cleaned stack trace and metadata
type CleanResult struct {
	Original  string       // Original stack trace
	Cleaned   string       // Cleaned stack trace
	Removed   int          // Number of repetitive blocks removed
	Frames    []StackFrame // Parsed stack frames
	ErrorInfo *ErrorInfo   // Extracted error information
}
