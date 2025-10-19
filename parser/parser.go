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

	// Performance optimization constants
	minStackTraceLength = 20
	maxLineLength       = 10000 // 10KB per line limit

	// Pattern matching constants
	minReactPatternMatches  = 4
	minSourcePatternMatches = 3
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
		if len(line) > maxLineLength {
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
// isStackFrameLine checks if a line contains a stack frame pattern
func isStackFrameLine(line string) bool {
	return framePattern.MatchString(line) || reactFramePattern.MatchString(line)
}

// countFrameOccurrences counts how many times each frame signature appears
func countFrameOccurrences(lines []string) map[string]int {
	frameCounts := make(map[string]int)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !isStackFrameLine(line) {
			continue
		}
		frameSignature := extractFrameSignature(line)
		frameCounts[frameSignature]++
	}
	return frameCounts
}

// buildCleanedLines creates cleaned lines, removing duplicate frames
func buildCleanedLines(lines []string) (cleanedLines []string, framesCollapsed int) {
	seenFrames := make(map[string]bool)

	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)

		if line == "" {
			// Preserve empty lines
			cleanedLines = append(cleanedLines, originalLine)
			continue
		}

		if !isStackFrameLine(line) {
			// Non-frame lines are preserved as-is
			cleanedLines = append(cleanedLines, originalLine)
			continue
		}

		frameSignature := extractFrameSignature(line)

		if seenFrames[frameSignature] {
			// This is a duplicate frame - skip it
			framesCollapsed++
			continue
		}

		// First occurrence of this frame - mark as seen and add normally
		seenFrames[frameSignature] = true
		cleanedLines = append(cleanedLines, originalLine)
	}

	return cleanedLines, framesCollapsed
}

// annotateDuplicateFrames adds annotations to frames that had duplicates
func annotateDuplicateFrames(cleanedLines []string, frameCounts map[string]int) {
	for i, line := range cleanedLines {
		lineTrimmed := strings.TrimSpace(line)
		if lineTrimmed == "" {
			continue
		}

		if isStackFrameLine(lineTrimmed) {
			frameSignature := extractFrameSignature(lineTrimmed)
			if count := frameCounts[frameSignature]; count > 1 {
				// This frame has duplicates - annotate it
				collapsedLine := fmt.Sprintf("%s // [x%d]", line, count)
				cleanedLines[i] = collapsedLine
			}
		}
	}
}

// extractErrorMessage extracts the error message from the first line of stack trace
func extractErrorMessage(lines []string) string {
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

// shouldIncludeSource checks if a source file should be included (not React internal)
func shouldIncludeSource(filename string) bool {
	return !strings.Contains(filename, "react-dom") &&
		!strings.Contains(filename, "ReactErrorUtils")
}

// extractSourceInfo attempts to extract source file and line information from a stack frame line
func extractSourceInfo(line string) string {
	// Try React console format first: "functionName @ file.js:123"
	if reactMatches := sourceFileReactPattern.FindStringSubmatch(line); len(reactMatches) >= minSourcePatternMatches {
		filename := strings.TrimSpace(reactMatches[1])
		if shouldIncludeSource(filename) {
			lineNum := reactMatches[2]
			return fmt.Sprintf("%s:%s", filename, lineNum)
		}
	}

	// Try alternative format
	if altMatches := sourceFileAltPattern.FindStringSubmatch(line); len(altMatches) >= minAltPatternMatches {
		filename := strings.TrimSpace(altMatches[1])
		if shouldIncludeSource(filename) {
			lineNum := altMatches[2]
			return fmt.Sprintf("%s:%s", filename, lineNum)
		}
	}

	// Try primary pattern for .js files (fallback)
	if jsMatches := sourceFilePattern.FindStringSubmatch(line); len(jsMatches) >= minJSPatternMatches {
		jsIndex := strings.LastIndex(line, ".js:")
		if jsIndex != -1 {
			start := strings.LastIndex(line[:jsIndex], "(")
			if start != -1 {
				filename := line[start+1 : jsIndex+3] // Include .js
				if shouldIncludeSource(filename) {
					return fmt.Sprintf("%s:%s", strings.TrimSpace(filename), jsMatches[1])
				}
			}
		}
	}

	return ""
}

// extractComponentInfo attempts to extract React component information from a stack frame line
func extractComponentInfo(line string) string {
	if matches := componentPattern.FindStringSubmatch(line); len(matches) >= minComponentPatternMatches {
		return matches[1]
	}
	return ""
}

// Only redundant stack frames are collapsed in-place - all other content including indentation and spacing is preserved exactly.
// Optimized to minimize string allocations and improve performance.
// Returns both the cleaned content and the exact count of frames collapsed.
func CleanStackTrace(content string) CleanResultPair {
	// Validate content before processing
	if !isValidContent(content) {
		return CleanResultPair{Content: content, Removed: 0}
	}

	if !IsStackTrace(content) {
		return CleanResultPair{Content: content, Removed: 0}
	}

	lines := strings.Split(content, "\n")
	frameCounts := countFrameOccurrences(lines)
	cleanedLines, framesCollapsed := buildCleanedLines(lines)
	annotateDuplicateFrames(cleanedLines, frameCounts)

	// Use strings.Builder for efficient string concatenation
	estimatedSize := len(content)
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
	return CleanResultPair{Content: result, Removed: framesCollapsed}
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
	if len(reactMatches) >= minReactPatternMatches {
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
	message := extractErrorMessage(lines)
	var stackFrames []string
	var source string
	var component string

	// Extract stack frames and look for React component info
	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to extract source information
		if source == "" {
			source = extractSourceInfo(line)
		}

		// Look for React component names
		if component == "" {
			component = extractComponentInfo(line)
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
