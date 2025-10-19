package models

import "time"

// ErrorInfo contains the essential information from a stack trace
type ErrorInfo struct {
	Stack     []string // Individual stack frames (24 bytes - pointer + len + cap)
	Message   string   // The main error message (8 bytes)
	Source    string   // Original source file/line if available (8 bytes)
	Component string   // React component name if applicable (8 bytes)
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
	Frames    []StackFrame // Parsed stack frames (24 bytes - pointer + len + cap)
	ErrorInfo *ErrorInfo   // Extracted error information (8 bytes)
	Original  string       // Original stack trace (8 bytes)
	Cleaned   string       // Cleaned stack trace (8 bytes)
	Removed   int          // Number of repetitive blocks removed (8 bytes)
}
