package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"com.github/rethunk-tech/tracetrim/internal/models"
)

const (
	// Minimum stack lines to consider content a stack trace
	minStackLinesForDetection = 2

	// Pattern matching constants
	minFunctionPatternMatches  = 4
	minAltPatternMatches       = 4
	minJSPatternMatches        = 3
	minComponentPatternMatches = 2

	// Estimation constants
	charsPerStackFrame = 50

	// Performance optimization constants
	minStackTraceLength = 20
	commentBufferExtra  = 50
)

// Pre-compiled regex patterns for better performance
var (
	// Stack trace detection patterns - enhanced for better edge case handling
	stackTracePatterns = []*regexp.Regexp{
		// JavaScript stack trace patterns - more precise to avoid false matches
		regexp.MustCompile(`\bat\s+[\w<>.()\s]+\s*\([^)]+\)`),          // "at functionName (file.js:123:45)" - allow more chars in function names
		regexp.MustCompile(`\b\w+\.(js|ts|jsx|tsx|mjs):\d+:\d+\b`),     // Support more file extensions
		regexp.MustCompile(`(?m)^Error:\s+.*\n\s+at\s+`),               // "Error: message\n    at" - multiline
		regexp.MustCompile(`\breact-dom\.development\.js`),             // React DOM development file
		regexp.MustCompile(`\bReactErrorUtils\.invokeGuardedCallback`), // Common React error pattern
		regexp.MustCompile(`\bUncaught\s+`),                            // "Uncaught Error:"
		regexp.MustCompile(`\bReferenceError:`),                        // "ReferenceError:"
		regexp.MustCompile(`\bTypeError:`),                             // "TypeError:"
		regexp.MustCompile(`\bSyntaxError:`),                           // "SyntaxError:"
		regexp.MustCompile(`\bEvalError:`),                             // "EvalError:"
		// React console output patterns
		regexp.MustCompile(`\b\w+\s+@\s+.+?:\d+\b`),                // "functionName @ file:line" - React console format (more permissive)
		regexp.MustCompile(`\b\w+\.(js|ts|jsx|tsx|mjs|cjs):\d+\b`), // File paths with line numbers (without column)
	}

	// Frame parsing patterns - enhanced for better edge case handling
	framePattern      = regexp.MustCompile(`(.+?)\s*\(([^:()]+):(\d+):(\d+)\)`)
	sourceFilePattern = regexp.MustCompile(`\.(js|ts|jsx|tsx|mjs|cjs):(\d+):(\d+)`)
	// React console format patterns
	reactFramePattern = regexp.MustCompile(`(.+?)\s*@\s*(.+?):(\d+)`)
	// Enhanced component patterns for React lifecycle methods
	componentPattern = regexp.MustCompile(`(\w+)\.(render|componentDidMount|componentDidUpdate|componentWillUnmount)\s*\(`)
	// Additional pattern for source file extraction with better path handling
	sourceFileAltPattern = regexp.MustCompile(`\(([^:()]+):(\d+):(\d+)\)`)
	// React console format for source file extraction
	sourceFileReactPattern = regexp.MustCompile(`@\s*(.+?):(\d+)`)
)

// IsStackTrace determines if the given content contains a JavaScript or React stack trace
// Optimized to avoid allocations for short content and improve performance
func IsStackTrace(content string) bool {
	// Validate input content first
	if !isValidContent(content) {
		return false
	}

	// Early return for obviously non-stack-trace content
	if len(content) < minStackTraceLength {
		return false
	}

	lines := strings.Split(content, "\n")
	stackLineCount := 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Check if this line matches any stack trace pattern
		for _, pattern := range stackTracePatterns {
			if pattern.MatchString(line) {
				stackLineCount++
				break
			}
		}

		// If we find multiple stack-like lines, it's likely a stack trace
		if stackLineCount >= minStackLinesForDetection {
			return true
		}
	}

	return false
}

// isValidContent validates that the content is safe to process
func isValidContent(content string) bool {
	// Check if content is valid UTF-8
	if !utf8.ValidString(content) {
		return false
	}

	// Check for null bytes (potential binary data)
	if strings.Contains(content, "\x00") {
		return false
	}

	// Check content length is reasonable (prevent memory exhaustion)
	if len(content) > 50*1024*1024 { // 50MB limit
		return false
	}

	// Check for extremely long lines that could cause issues
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if len(line) > 10000 { // 10KB per line limit
			return false
		}
	}

	return true
}

// CleanResultPair contains both the cleaned content and detailed statistics
type CleanResultPair struct {
	Content string
	Removed int
}

// CleanStackTrace removes repetitive stack trace blocks while preserving all original formatting.
// Only redundant stack frames are removed - all other content including indentation and spacing is preserved exactly.
// Optimized to minimize string allocations and improve performance.
// Returns both the cleaned content and the exact count of frames removed.
func CleanStackTrace(content string) CleanResultPair {
	// Validate content before processing
	if !isValidContent(content) {
		return CleanResultPair{Content: content, Removed: 0}
	}

	if !IsStackTrace(content) {
		return CleanResultPair{Content: content, Removed: 0}
	}

	lines := strings.Split(content, "\n")
	var cleanedLines []string
	seenFrames := make(map[string]bool)
	var framesRemoved int

	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" {
			// Preserve empty lines
			cleanedLines = append(cleanedLines, originalLine)
			continue
		}

		// Check if this is a stack frame line (contains file:line:column pattern or React console format)
		if framePattern.MatchString(line) || reactFramePattern.MatchString(line) {
			// Extract the frame signature (function + file + line)
			frameSignature := extractFrameSignature(line)

			if seenFrames[frameSignature] {
				framesRemoved++
				continue // Skip this duplicate frame
			} else {
				seenFrames[frameSignature] = true
				cleanedLines = append(cleanedLines, originalLine)
			}
		} else {
			// Non-frame line (error message, etc.) - always include
			cleanedLines = append(cleanedLines, originalLine)
		}
	}

	// Use strings.Builder for efficient string concatenation
	// Pre-calculate approximate size to minimize reallocations
	estimatedSize := len(content) // Start with original size
	var builder strings.Builder
	builder.Grow(estimatedSize)

	// Join cleaned lines
	for i, line := range cleanedLines {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(line)
	}

	result := builder.String()

	// If we removed duplicates, add a note about it
	if framesRemoved > 0 {
		// Reset builder and rebuild with the note
		builder.Reset()
		builder.Grow(estimatedSize + commentBufferExtra) // Extra space for the comment
		builder.WriteString(fmt.Sprintf("// Removed %d repetitive stack frame(s)\n", framesRemoved))
		builder.WriteString(result)
		result = builder.String()
	}

	return CleanResultPair{
		Content: result,
		Removed: framesRemoved,
	}
}

// extractFrameSignature creates a unique signature for a stack frame to detect duplicates
// This function is now optimized to avoid redundant regex compilation and string processing
func extractFrameSignature(line string) string {
	// Try standard format first: "at functionName (file.js:123:45)"
	matches := framePattern.FindStringSubmatch(line)
	if len(matches) >= minFunctionPatternMatches {
		// Use helper function for consistent React internal function handling
		functionName := strings.TrimSpace(matches[1])
		functionName = strings.TrimPrefix(functionName, "at ")
		return extractFrameSignatureForStandardFormat(functionName, matches[2], matches[3])
	}

	// Try React console format: "functionName @ file.js:123"
	reactMatches := reactFramePattern.FindStringSubmatch(line)
	if len(reactMatches) >= 4 {
		functionName := strings.TrimSpace(reactMatches[1])
		fileName := reactMatches[2]
		lineNumber := reactMatches[3]

		// For React internal functions, use function + file (ignoring line) to detect duplicates
		// This helps collapse repeated React internal calls that happen on different lines
		if isReactInternalFunction(functionName, fileName) {
			return fmt.Sprintf("%s|%s", functionName, fileName)
		}

		return fmt.Sprintf("%s|%s|%s", functionName, fileName, lineNumber)
	}

	return line // Fallback to entire line if parsing fails
}

// isReactInternalFunction determines if a function is a React internal function
// that should have its line numbers ignored for duplicate detection
func isReactInternalFunction(functionName, fileName string) bool {
	// React DOM development files contain many internal functions that are called repeatedly
	if strings.Contains(fileName, "react-dom") {
		// Common React internal functions that appear repeatedly in stack traces
		reactInternalFunctions := []string{
			"recursivelyTraverseAndDoubleInvokeEffectsInDEV",
			"recursivelyTraversePassiveMountEffects",
			"commitPassiveMountOnFiber",
			"recursivelyTraverseReconnectPassiveEffects",
			"recursivelyTraverseDisconnectPassiveEffects",
			"recursivelyTraversePassiveUnmountEffects",
			"commitPassiveUnmountOnFiber",
			"ReactErrorUtils.invokeGuardedCallback",
			"ReactCompositeComponent._renderValidatedComponent",
			"react_stack_bottom_frame",
		}

		functionNameLower := strings.ToLower(functionName)
		for _, internalFunc := range reactInternalFunctions {
			if strings.Contains(functionNameLower, strings.ToLower(internalFunc)) {
				return true
			}
		}
	}

	return false
}

// extractFrameSignatureForStandardFormat handles standard format frame signatures
func extractFrameSignatureForStandardFormat(functionName, fileName, lineNumber string) string {
	// For React internal functions, use function + file (ignoring line) to detect duplicates
	if isReactInternalFunction(functionName, fileName) {
		return fmt.Sprintf("%s|%s", functionName, fileName)
	}

	return fmt.Sprintf("%s|%s|%s", functionName, fileName, lineNumber)
}

// ExtractErrorInfo extracts structured information from a stack trace for analysis.
// This function only parses and analyzes the content - it does not modify the original clipboard content.
func ExtractErrorInfo(content string) *models.ErrorInfo {
	if !IsStackTrace(content) {
		return nil
	}

	lines := strings.Split(content, "\n")
	var message string
	var stackFrames []string
	var source string
	var component string

	// Extract error message (usually the first line)
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		message = firstLine
	}

	// Extract stack frames and look for React component info
	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract file and line info for source - prioritize user code over React internals
		// Try React console format first: "functionName @ file.js:123"
		if reactMatches := sourceFileReactPattern.FindStringSubmatch(line); len(reactMatches) >= 3 {
			filename := strings.TrimSpace(reactMatches[1])

			// Only set source if this is not React internal code
			if !strings.Contains(filename, "react-dom") &&
				!strings.Contains(filename, "ReactErrorUtils") {
				lineNum := reactMatches[2]
				source = fmt.Sprintf("%s:%s", filename, lineNum)
			}
		} else if altMatches := sourceFileAltPattern.FindStringSubmatch(line); len(altMatches) >= minAltPatternMatches {
			filename := strings.TrimSpace(altMatches[1])

			// Only set source if this is not React internal code
			if !strings.Contains(filename, "react-dom") &&
				!strings.Contains(filename, "ReactErrorUtils") {
				lineNum := altMatches[2]
				source = fmt.Sprintf("%s:%s", filename, lineNum)
			}
		}

		// Also try primary pattern for .js files (fallback)
		if source == "" {
			if jsMatches := sourceFilePattern.FindStringSubmatch(line); len(jsMatches) >= minJSPatternMatches {
				jsIndex := strings.LastIndex(line, ".js:")
				if jsIndex != -1 {
					start := strings.LastIndex(line[:jsIndex], "(")
					if start != -1 {
						filename := line[start+1 : jsIndex+3] // Include .js

						// Only set if not React internal code
						if !strings.Contains(filename, "react-dom") &&
							!strings.Contains(filename, "ReactErrorUtils") {
							source = fmt.Sprintf("%s:%s", strings.TrimSpace(filename), jsMatches[1])
						}
					}
				}
			}
		}

		// Look for React component names using enhanced pattern (includes lifecycle methods)
		if matches := componentPattern.FindStringSubmatch(line); len(matches) >= minComponentPatternMatches {
			component = matches[1]
		}

		stackFrames = append(stackFrames, originalLine)
	}

	return &models.ErrorInfo{
		Message:   message,
		Stack:     stackFrames,
		Source:    source,
		Component: component,
	}
}

// CleanResult provides detailed information about the cleaning process
func CleanResult(content string) models.CleanResult {
	original := content

	// Use the enhanced CleanStackTrace that returns accurate frame removal count
	cleanResult := CleanStackTrace(content)
	cleaned := cleanResult.Content
	removed := cleanResult.Removed

	// Calculate bytes saved
	bytesSaved := len(original) - len(cleaned)

	// Count lines before and after
	linesBefore := strings.Count(original, "\n") + 1 // +1 for the last line if no trailing newline
	linesAfter := strings.Count(cleaned, "\n") + 1

	var frames []models.StackFrame
	errorInfo := ExtractErrorInfo(content)

	return models.CleanResult{
		Original:    original,
		Cleaned:     cleaned,
		Removed:     removed,
		BytesSaved:  bytesSaved,
		LinesBefore: linesBefore,
		LinesAfter:  linesAfter,
		Frames:      frames,
		ErrorInfo:   errorInfo,
	}
}
