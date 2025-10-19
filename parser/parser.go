package parser

import (
	"fmt"
	"regexp"
	"strings"

	"com.github/rethunk-tech/no-reaction/internal/models"
)

const (
	// Minimum stack lines to consider content a stack trace
	minStackLinesForDetection = 2

	// Pattern matching constants
	minFunctionPatternMatches = 4
	minAltPatternMatches      = 4
	minJSPatternMatches       = 3

	// Estimation constants
	charsPerStackFrame = 50
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
	}

	// Frame parsing patterns - enhanced for better edge case handling
	framePattern      = regexp.MustCompile(`(.+?)\s*\(([^:()]+):(\d+):(\d+)\)`)
	sourceFilePattern = regexp.MustCompile(`\.(js|ts|jsx|tsx|mjs):(\d+):(\d+)`)
	// Enhanced component patterns for React lifecycle methods
	componentPattern = regexp.MustCompile(`(\w+)\.(render|componentDidMount|componentDidUpdate|componentWillUnmount)\s*\(`)
	// Additional pattern for source file extraction with better path handling
	sourceFileAltPattern = regexp.MustCompile(`\(([^:()]+):(\d+):(\d+)\)`)
)

// IsStackTrace determines if the given content contains a JavaScript or React stack trace
func IsStackTrace(content string) bool {
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

// CleanStackTrace removes repetitive stack trace blocks while preserving all original formatting.
// Only redundant stack frames are removed - all other content including indentation and spacing is preserved exactly.
func CleanStackTrace(content string) string {
	if !IsStackTrace(content) {
		return content
	}

	lines := strings.Split(content, "\n")
	var cleanedLines []string
	seenFrames := make(map[string]bool)
	var consecutiveDuplicates int

	for _, line := range lines {
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" {
			// Preserve empty lines
			cleanedLines = append(cleanedLines, originalLine)
			continue
		}

		// Check if this is a stack frame line (contains file:line:column pattern)
		if framePattern.MatchString(line) {
			// Extract the frame signature (function + file + line)
			frameSignature := extractFrameSignature(line)

			if seenFrames[frameSignature] {
				consecutiveDuplicates++
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
	var builder strings.Builder

	// Join cleaned lines
	for i, line := range cleanedLines {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(line)
	}

	result := builder.String()

	// If we removed duplicates, add a note about it
	if consecutiveDuplicates > 0 {
		// Reset builder and rebuild with the note
		builder.Reset()
		builder.WriteString(fmt.Sprintf("// Removed %d repetitive stack frame(s)\n", consecutiveDuplicates))
		builder.WriteString(result)
		result = builder.String()
	}

	return result
}

// extractFrameSignature creates a unique signature for a stack frame to detect duplicates
// This function is now optimized to avoid redundant regex compilation and string processing
func extractFrameSignature(line string) string {
	matches := framePattern.FindStringSubmatch(line)
	if len(matches) >= minFunctionPatternMatches {
		// Return function + file + line as unique signature (trim "at " prefix if present)
		functionName := strings.TrimSpace(matches[1])
		if strings.HasPrefix(functionName, "at ") {
			functionName = strings.TrimPrefix(functionName, "at ")
		}
		return fmt.Sprintf("%s|%s|%s", functionName, matches[2], matches[3])
	}

	return line // Fallback to entire line if parsing fails
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
		// Try alternative pattern first (more reliable for all file types)
		if altMatches := sourceFileAltPattern.FindStringSubmatch(line); len(altMatches) >= minAltPatternMatches {
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
		if matches := componentPattern.FindStringSubmatch(line); len(matches) >= 2 {
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
	cleaned := CleanStackTrace(content)
	removed := 0

	// Simple heuristic: if cleaned version is shorter, we removed something
	if len(cleaned) < len(original) {
		removed = (len(original) - len(cleaned)) / charsPerStackFrame // Rough estimate
	}

	var frames []models.StackFrame
	errorInfo := ExtractErrorInfo(content)

	return models.CleanResult{
		Original:  original,
		Cleaned:   cleaned,
		Removed:   removed,
		Frames:    frames,
		ErrorInfo: errorInfo,
	}
}
