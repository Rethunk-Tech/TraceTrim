package parser

import (
	"fmt"
	"regexp"
	"strings"

	"com.github/rethunk-tech/no-reaction/internal/models"
)

// IsStackTrace determines if the given content contains a JavaScript or React stack trace
func IsStackTrace(content string) bool {
	// Look for common stack trace patterns
	stackTracePatterns := []string{
		// JavaScript stack trace patterns
		`at\s+[\w<>.\s]+\([^)]+\)`, // "at functionName (file.js:123:45)"
		`\w+\.js:\d+:\d+`,          // "file.js:123:45"
		`Error:\s+.*\n\s+at\s+`,    // "Error: message\n    at"
		// React specific patterns
		`react-dom\.development\.js`,             // React DOM development file
		`ReactErrorUtils\.invokeGuardedCallback`, // Common React error pattern
		// Generic error patterns
		`Uncaught\s+`,     // "Uncaught Error:"
		`ReferenceError:`, // "ReferenceError:"
		`TypeError:`,      // "TypeError:"
	}

	lines := strings.Split(content, "\n")
	stackLineCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line matches any stack trace pattern
		for _, pattern := range stackTracePatterns {
			matched, _ := regexp.MatchString(pattern, line)
			if matched {
				stackLineCount++
				break
			}
		}

		// If we find multiple stack-like lines, it's likely a stack trace
		if stackLineCount >= 3 {
			return true
		}
	}

	return false
}

// CleanStackTrace removes repetitive blocks from stack traces while preserving essential information
func CleanStackTrace(content string) string {
	if !IsStackTrace(content) {
		return content
	}

	lines := strings.Split(content, "\n")
	var cleanedLines []string
	seenFrames := make(map[string]bool)
	var consecutiveDuplicates int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this is a stack frame line (contains file:line:column pattern)
		framePattern := `(.+?)\s*\(([^:]+):(\d+):(\d+)\)`
		frameRegex := regexp.MustCompile(framePattern)

		if frameRegex.MatchString(line) {
			// Extract the frame signature (function + file + line)
			frameSignature := extractFrameSignature(line)

			if seenFrames[frameSignature] {
				consecutiveDuplicates++
				continue // Skip this duplicate frame
			} else {
				seenFrames[frameSignature] = true
				cleanedLines = append(cleanedLines, line)
			}
		} else {
			// Non-frame line (error message, etc.) - always include
			cleanedLines = append(cleanedLines, line)
		}
	}

	result := strings.Join(cleanedLines, "\n")

	// If we removed duplicates, add a note about it
	if consecutiveDuplicates > 0 {
		result = fmt.Sprintf("// Removed %d repetitive stack frame(s)\n%s", consecutiveDuplicates, result)
	}

	return result
}

// extractFrameSignature creates a unique signature for a stack frame to detect duplicates
func extractFrameSignature(line string) string {
	// Extract function name, file, and line number
	framePattern := `(.+?)\s*\(([^:]+):(\d+):(\d+)\)`
	frameRegex := regexp.MustCompile(framePattern)

	matches := frameRegex.FindStringSubmatch(line)
	if len(matches) >= 4 {
		// Return function + file + line as unique signature
		return fmt.Sprintf("%s|%s|%s", matches[1], matches[2], matches[3])
	}

	return line // Fallback to entire line if parsing fails
}

// ExtractErrorInfo extracts structured information from a stack trace
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
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract file and line info for source
		filePattern := `([^:]+):(\d+):(\d+)`
		fileRegex := regexp.MustCompile(filePattern)

		if matches := fileRegex.FindStringSubmatch(line); len(matches) >= 3 {
			source = fmt.Sprintf("%s:%s", matches[1], matches[2])
		}

		// Look for React component names
		componentPattern := `in\s+(\w+)`
		componentRegex := regexp.MustCompile(componentPattern)

		if matches := componentRegex.FindStringSubmatch(line); len(matches) >= 2 {
			component = matches[1]
		}

		stackFrames = append(stackFrames, line)
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
		removed = (len(original) - len(cleaned)) / 50 // Rough estimate
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
