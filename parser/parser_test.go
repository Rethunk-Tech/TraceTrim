package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestIsStackTrace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name: "React error with repetitive frames",
			input: "Error: Something went wrong\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)",
			expected: true,
		},
		{
			name:     "JavaScript error",
			input:    "ReferenceError: x is not defined\n    at eval (eval at <anonymous> (test.js:1:1))",
			expected: true,
		},
		{
			name: "Multiple stack frames",
			input: "TypeError: Cannot read property 'map' of undefined\n" +
				"    at Component.render (app.js:15:20)\n" +
				"    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)",
			expected: true,
		},
		{
			name:     "Plain text (not a stack trace)",
			input:    "This is just some regular text with no stack trace",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Single line that looks like stack frame",
			input:    "    at someFunction (file.js:123:45)",
			expected: false,
		},
		{
			name: "TypeScript file extension",
			input: "Error: Type error\n" +
				"    at myFunction (app.ts:15:10)\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)",
			expected: true,
		},
		{
			name: "React JSX file extension",
			input: "Error: Component error\n" +
				"    at MyComponent.render (Component.jsx:25:5)\n" +
				"    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)",
			expected: true,
		},
		{
			name: "React console output format",
			input: "useAuth.useEffect @ S:\\Projects\\com.github\\PeleOs-LLC\\ROK-UI-v2\\src\\lib\\hooks\\useAuth.ts:47\n" +
				"react_stack_bottom_frame @ react-dom-client.development.js:23669\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStackTrace(tt.input)
			if result != tt.expected {
				t.Errorf("IsStackTrace(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// runStringTests runs table-driven tests for functions that take a string and return a string
func runStringTests(t *testing.T, tests []struct {
	name     string
	input    string
	expected string
}, testFunc func(string) string,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testFunc(tt.input)
			if result != tt.expected {
				t.Errorf("%s(%q) =\n%q, want\n%q", t.Name(), tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanStackTrace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Remove repetitive React frames",
			input: `Error: Failed to render
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`,
			expected: `// Removed 2 repetitive stack frame(s)
Error: Failed to render
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`,
		},
		{
			name: "No duplicates to remove",
			input: `Error: Something failed
    at function1 (file1.js:10:5)
    at function2 (file2.js:20:10)
    at function3 (file3.js:30:15)`,
			expected: `Error: Something failed
    at function1 (file1.js:10:5)
    at function2 (file2.js:20:10)
    at function3 (file3.js:30:15)`,
		},
		{
			name:     "Not a stack trace",
			input:    "This is just regular text",
			expected: "This is just regular text",
		},
		{
			name: "Mixed content with some duplicates",
			input: `TypeError: Cannot read property 'name' of undefined
    at UserProfile.render (UserProfile.js:45:12)
    at UserProfile.render (UserProfile.js:45:12)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
    at UserProfile.render (UserProfile.js:45:12)`,
			expected: `// Removed 2 repetitive stack frame(s)
TypeError: Cannot read property 'name' of undefined
    at UserProfile.render (UserProfile.js:45:12)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`,
		},
		{
			name: "React console format with duplicates",
			input: `useAuth.useEffect @ S:\Projects\com.github\PeleOs-LLC\ROK-UI-v2\src\lib\hooks\useAuth.ts:47
react_stack_bottom_frame @ react-dom-client.development.js:23669
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
react_stack_bottom_frame @ react-dom-client.development.js:23669
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`,
			expected: `// Removed 2 repetitive stack frame(s)
useAuth.useEffect @ S:\Projects\com.github\PeleOs-LLC\ROK-UI-v2\src\lib\hooks\useAuth.ts:47
react_stack_bottom_frame @ react-dom-client.development.js:23669
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanStackTrace(tt.input)

			// For backward compatibility, check if the content matches expected
			// The new function returns a CleanResultPair with Content field
			content := result.Content
			if content != tt.expected {
				t.Errorf("CleanStackTrace(%q) =\n%q, want\n%q", tt.input, content, tt.expected)
			}
		})
	}
}

func TestExtractErrorInfo(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedMsg  string
		expectedSrc  string
		expectedComp string
	}{
		{
			name: "React component error",
			input: "Error: Objects are not valid as a React child\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)\n" +
				"    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)\n" +
				"    at MyComponent.render (MyComponent.js:25:10)",
			expectedMsg:  "Error: Objects are not valid as a React child",
			expectedSrc:  "MyComponent.js:25",
			expectedComp: "MyComponent",
		},
		{
			name: "React component lifecycle method",
			input: "Warning: Component did update\n" +
				"    at MyComponent.componentDidUpdate (MyComponent.js:45:8)\n" +
				"    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)",
			expectedMsg:  "Warning: Component did update",
			expectedSrc:  "MyComponent.js:45",
			expectedComp: "MyComponent",
		},
		{
			name:         "JavaScript error",
			input:        "ReferenceError: x is not defined\n    at eval (eval at <anonymous> (script.js:1:1))",
			expectedMsg:  "ReferenceError: x is not defined",
			expectedSrc:  "script.js:1",
			expectedComp: "",
		},
		{
			name:         "Not a stack trace",
			input:        "This is not a stack trace",
			expectedMsg:  "",
			expectedSrc:  "",
			expectedComp: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractErrorInfo(tt.input)

			if tt.expectedMsg == "" {
				if result != nil {
					t.Errorf("ExtractErrorInfo(%q) expected nil, got %+v", tt.input, result)
				}
				return
			}

			if result == nil {
				t.Errorf("ExtractErrorInfo(%q) expected non-nil result", tt.input)
				return
			}

			if result.Message != tt.expectedMsg {
				t.Errorf("ExtractErrorInfo(%q).Message = %q, want %q", tt.input, result.Message, tt.expectedMsg)
			}

			if result.Source != tt.expectedSrc {
				t.Errorf("ExtractErrorInfo(%q).Source = %q, want %q", tt.input, result.Source, tt.expectedSrc)
			}

			if result.Component != tt.expectedComp {
				t.Errorf("ExtractErrorInfo(%q).Component = %q, want %q", tt.input, result.Component, tt.expectedComp)
			}
		})
	}
}

func TestExtractFrameSignature(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard frame format",
			input:    "    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)",
			expected: "ReactErrorUtils.invokeGuardedCallback|react-dom.development.js",
		},
		{
			name:     "Different format",
			input:    "    at Component.render (MyComponent.js:25:10)",
			expected: "Component.render|MyComponent.js|25",
		},
		{
			name:     "Anonymous function",
			input:    "    at <anonymous> (script.js:5:12)",
			expected: "<anonymous>|script.js|5",
		},
		{
			name:     "No parentheses",
			input:    "    at someFunction file.js:10:5",
			expected: "    at someFunction file.js:10:5", // Should return original if parsing fails
		},
		{
			name:     "React console format",
			input:    "react_stack_bottom_frame @ react-dom-client.development.js:23669",
			expected: "react_stack_bottom_frame|react-dom-client.development.js",
		},
		{
			name:     "React console format with Windows path",
			input:    "useAuth.useEffect @ S:\\Projects\\com.github\\PeleOs-LLC\\ROK-UI-v2\\src\\lib\\hooks\\useAuth.ts:47",
			expected: "useAuth.useEffect|S:\\Projects\\com.github\\PeleOs-LLC\\ROK-UI-v2\\src\\lib\\hooks\\useAuth.ts|47",
		},
		{
			name:     "React internal function (should ignore line number)",
			input:    "commitPassiveMountOnFiber @ react-dom-client.development.js:14380",
			expected: "commitPassiveMountOnFiber|react-dom-client.development.js",
		},
		{
			name:     "React internal function different line (should be same signature)",
			input:    "commitPassiveMountOnFiber @ react-dom-client.development.js:14514",
			expected: "commitPassiveMountOnFiber|react-dom-client.development.js",
		},
	}

	runStringTests(t, tests, extractFrameSignature)
}

func TestCleanResult(t *testing.T) {
	input := `Error: Test error
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`

	result := CleanResult(input)

	// Check that we have a result
	if result.Original != input {
		t.Errorf("CleanResult.Original = %q, want %q", result.Original, input)
	}

	// Check that cleaned version is different (duplicates removed)
	if result.Cleaned == input {
		t.Error("CleanResult.Cleaned should be different from original when duplicates are removed")
	}

	// Check that we detected some removals
	if result.Removed <= 0 {
		t.Error("CleanResult.Removed should be > 0 when duplicates are removed")
	}

	// Check that bytes saved is calculated correctly
	expectedBytesSaved := len(input) - len(result.Cleaned)
	if result.BytesSaved != expectedBytesSaved {
		t.Errorf("CleanResult.BytesSaved = %d, want %d", result.BytesSaved, expectedBytesSaved)
	}

	// Check that line counts are correct
	expectedLinesBefore := strings.Count(input, "\n") + 1
	expectedLinesAfter := strings.Count(result.Cleaned, "\n") + 1
	if result.LinesBefore != expectedLinesBefore {
		t.Errorf("CleanResult.LinesBefore = %d, want %d", result.LinesBefore, expectedLinesBefore)
	}
	if result.LinesAfter != expectedLinesAfter {
		t.Errorf("CleanResult.LinesAfter = %d, want %d", result.LinesAfter, expectedLinesAfter)
	}

	// Check that error info was extracted
	if result.ErrorInfo == nil {
		t.Error("CleanResult.ErrorInfo should not be nil for valid stack trace")
	} else if result.ErrorInfo.Message != "Error: Test error" {
		t.Errorf("CleanResult.ErrorInfo.Message = %q, want %q", result.ErrorInfo.Message, "Error: Test error")
	}
}

func TestCleanStackTraceReturnsAccurateCount(t *testing.T) {
	input := `Error: Test error
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`

	result := CleanStackTrace(input)

	// Check that we detected exactly 2 removals (3 duplicates, but we keep the first one)
	if result.Removed != 2 {
		t.Errorf("CleanStackTrace.Removed = %d, want 2", result.Removed)
	}

	// Check that cleaned content is different from original
	if result.Content == input {
		t.Error("CleanStackTrace.Content should be different from original when duplicates are removed")
	}

	// Check that the comment mentions the correct number
	if !strings.Contains(result.Content, "Removed 2 repetitive stack frame(s)") {
		t.Error("CleanStackTrace.Content should contain comment with correct removal count")
	}
}

func TestCleanStackTraceNoDuplicates(t *testing.T) {
	input := `Error: Test error
    at function1 (file1.js:10:5)
    at function2 (file2.js:20:10)
    at function3 (file3.js:30:15)`

	result := CleanStackTrace(input)

	// Check that no frames were removed
	if result.Removed != 0 {
		t.Errorf("CleanStackTrace.Removed = %d, want 0", result.Removed)
	}

	// Check that content is unchanged
	if result.Content != input {
		t.Error("CleanStackTrace.Content should be unchanged when no duplicates exist")
	}
}

func TestCleanResultStatistics(t *testing.T) {
	input := `Error: Test error
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`

	result := CleanResult(input)

	// Test bytes saved calculation
	if result.BytesSaved != len(input)-len(result.Cleaned) {
		t.Errorf("BytesSaved calculation is incorrect")
	}

	// Test lines count
	expectedLinesBefore := 3 // Error line + 2 stack frame lines
	if result.LinesBefore != expectedLinesBefore {
		t.Errorf("LinesBefore = %d, want %d", result.LinesBefore, expectedLinesBefore)
	}

	// After cleaning, should have 3 lines (Error + comment + 1 unique frame)
	expectedLinesAfter := 3
	if result.LinesAfter != expectedLinesAfter {
		t.Errorf("LinesAfter = %d, want %d", result.LinesAfter, expectedLinesAfter)
	}
}

// TestReactConsoleFormatDebug helps debug the React console format processing
func TestReactConsoleFormatDebug(t *testing.T) {
	input := `useAuth.useEffect @ S:\Projects\com.github\PeleOs-LLC\ROK-UI-v2\src\lib\hooks\useAuth.ts:47
react_stack_bottom_frame @ react-dom-client.development.js:23669
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
react_stack_bottom_frame @ react-dom-client.development.js:23669
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)`

	lines := strings.Split(input, "\n")
	frameSignatures := make(map[string][]string)

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		sig := extractFrameSignature(line)
		frameSignatures[sig] = append(frameSignatures[sig], fmt.Sprintf("Line %d: %s", i, line))
	}

	// Print debug info
	for sig, occurrences := range frameSignatures {
		t.Logf("Signature: %s, Occurrences: %d", sig, len(occurrences))
		for _, occurrence := range occurrences {
			t.Logf("  %s", occurrence)
		}
	}

	result := CleanStackTrace(input)
	t.Logf("Input: %q", input)
	t.Logf("Output: %q", result.Content)
	t.Logf("Removed: %d", result.Removed)
}

// Benchmark tests
func BenchmarkIsStackTrace(b *testing.B) {
	input := `Error: Something went wrong
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`

	for i := 0; i < b.N; i++ {
		IsStackTrace(input)
	}
}

func BenchmarkCleanStackTrace(b *testing.B) {
	input := `Error: Something went wrong
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`

	for i := 0; i < b.N; i++ {
		CleanStackTrace(input)
	}
}

func BenchmarkExtractErrorInfo(b *testing.B) {
	input := `Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)
    at MyComponent.render (MyComponent.js:25:10)`

	for i := 0; i < b.N; i++ {
		ExtractErrorInfo(input)
	}
}
