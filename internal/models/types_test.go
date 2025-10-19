package models

import (
	"strings"
	"testing"
	"time"
)

func TestClipboardContent(t *testing.T) {
	// Test creating clipboard content
	content := ClipboardContent{
		Content:   "test clipboard content",
		Timestamp: time.Now(),
		Format:    "text/plain",
	}

	if content.Content != "test clipboard content" {
		t.Errorf("Expected content 'test clipboard content', got %q", content.Content)
	}

	if content.Format != "text/plain" {
		t.Errorf("Expected format 'text/plain', got %q", content.Format)
	}

	if content.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	// Test with different format
	content.Format = "text/html"
	if content.Format != "text/html" {
		t.Errorf("Expected format 'text/html', got %q", content.Format)
	}
}

func TestErrorInfo(t *testing.T) {
	// Test creating error info
	errorInfo := ErrorInfo{
		Stack:     []string{"frame1", "frame2", "frame3"},
		Message:   "Test error message",
		Source:    "app.js:25",
		Component: "TestComponent",
	}

	if len(errorInfo.Stack) != 3 {
		t.Errorf("Expected 3 stack frames, got %d", len(errorInfo.Stack))
	}

	if errorInfo.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got %q", errorInfo.Message)
	}

	if errorInfo.Source != "app.js:25" {
		t.Errorf("Expected source 'app.js:25', got %q", errorInfo.Source)
	}

	if errorInfo.Component != "TestComponent" {
		t.Errorf("Expected component 'TestComponent', got %q", errorInfo.Component)
	}

	// Test with empty fields
	emptyErrorInfo := ErrorInfo{}
	if len(emptyErrorInfo.Stack) != 0 {
		t.Errorf("Expected empty stack, got %d frames", len(emptyErrorInfo.Stack))
	}
	if emptyErrorInfo.Message != "" {
		t.Errorf("Expected empty message, got %q", emptyErrorInfo.Message)
	}
}

func TestStackFrame(t *testing.T) {
	// Test creating stack frame
	frame := StackFrame{
		Function: "testFunction",
		File:     "test.js",
		Line:     42,
		Column:   15,
	}

	if frame.Function != "testFunction" {
		t.Errorf("Expected function 'testFunction', got %q", frame.Function)
	}

	if frame.File != "test.js" {
		t.Errorf("Expected file 'test.js', got %q", frame.File)
	}

	if frame.Line != 42 {
		t.Errorf("Expected line 42, got %d", frame.Line)
	}

	if frame.Column != 15 {
		t.Errorf("Expected column 15, got %d", frame.Column)
	}

	// Test with zero values
	zeroFrame := StackFrame{}
	if zeroFrame.Line != 0 {
		t.Errorf("Expected line 0, got %d", zeroFrame.Line)
	}
	if zeroFrame.Column != 0 {
		t.Errorf("Expected column 0, got %d", zeroFrame.Column)
	}
}

func TestCleanResult(t *testing.T) {
	// Test creating clean result
	cleanResult := CleanResult{
		Frames: []StackFrame{
			{Function: "func1", File: "file1.js", Line: 10, Column: 5},
			{Function: "func2", File: "file2.js", Line: 20, Column: 10},
		},
		ErrorInfo: &ErrorInfo{
			Message: "Test error",
			Source:  "app.js:1",
		},
		Original:    "original stack trace",
		Cleaned:     "cleaned stack trace",
		Removed:     5,
		BytesSaved:  100,
		LinesBefore: 10,
		LinesAfter:  8,
	}

	if len(cleanResult.Frames) != 2 {
		t.Errorf("Expected 2 frames, got %d", len(cleanResult.Frames))
	}

	if cleanResult.ErrorInfo == nil {
		t.Error("Expected error info to be set")
	} else if cleanResult.ErrorInfo.Message != "Test error" {
		t.Errorf("Expected error message 'Test error', got %q", cleanResult.ErrorInfo.Message)
	}

	if cleanResult.Original != "original stack trace" {
		t.Errorf("Expected original 'original stack trace', got %q", cleanResult.Original)
	}

	if cleanResult.Cleaned != "cleaned stack trace" {
		t.Errorf("Expected cleaned 'cleaned stack trace', got %q", cleanResult.Cleaned)
	}

	if cleanResult.Removed != 5 {
		t.Errorf("Expected removed 5, got %d", cleanResult.Removed)
	}

	if cleanResult.BytesSaved != 100 {
		t.Errorf("Expected bytes saved 100, got %d", cleanResult.BytesSaved)
	}

	if cleanResult.LinesBefore != 10 {
		t.Errorf("Expected lines before 10, got %d", cleanResult.LinesBefore)
	}

	if cleanResult.LinesAfter != 8 {
		t.Errorf("Expected lines after 8, got %d", cleanResult.LinesAfter)
	}

	// Test with nil ErrorInfo
	cleanResultNilError := CleanResult{
		ErrorInfo: nil,
	}

	if cleanResultNilError.ErrorInfo != nil {
		t.Error("Expected error info to be nil")
	}
}

func TestModelFieldAlignment(t *testing.T) {
	// This test ensures that the struct fields are properly aligned
	// and can be created and accessed correctly

	// Test multiple instances
	contents := []ClipboardContent{
		{Content: "content1", Format: "text/plain"},
		{Content: "content2", Format: "text/html"},
		{Content: "content3", Format: "application/json"},
	}

	for i, content := range contents {
		expectedContent := "content" + string(rune('1'+i))
		if content.Content != expectedContent {
			t.Errorf("Expected content %q, got %q", expectedContent, content.Content)
		}
	}

	// Test error info stack manipulation
	errorInfo := ErrorInfo{
		Stack: []string{"frame1", "frame2"},
	}

	errorInfo.Stack = append(errorInfo.Stack, "frame3")
	if len(errorInfo.Stack) != 3 {
		t.Errorf("Expected 3 stack frames after append, got %d", len(errorInfo.Stack))
	}

	// Test clean result calculations
	original := "line1\nline2\nline3"
	cleaned := "line1\nline3"
	cleanResult := CleanResult{
		Original:   original,
		Cleaned:    cleaned,
		BytesSaved: len(original) - len(cleaned), // Calculate actual difference
	}

	if cleanResult.BytesSaved != len(cleanResult.Original)-len(cleanResult.Cleaned) {
		t.Errorf("BytesSaved should equal original length minus cleaned length: %d - %d = %d, but got %d",
			len(cleanResult.Original), len(cleanResult.Cleaned),
			len(cleanResult.Original)-len(cleanResult.Cleaned), cleanResult.BytesSaved)
	}

	if cleanResult.LinesBefore != strings.Count(cleanResult.Original, "\n")+1 {
		t.Error("LinesBefore should be correct line count")
	}

	if cleanResult.LinesAfter != strings.Count(cleanResult.Cleaned, "\n")+1 {
		t.Error("LinesAfter should be correct line count")
	}
}

func TestTimestampOperations(t *testing.T) {
	// Test timestamp operations on ClipboardContent
	now := time.Now()
	content := ClipboardContent{
		Timestamp: now,
	}

	if !content.Timestamp.Equal(now) {
		t.Error("Timestamp should be set correctly")
	}

	// Test that timestamps are properly comparable
	content2 := ClipboardContent{
		Timestamp: now.Add(time.Second),
	}

	if content.Timestamp.Equal(content2.Timestamp) {
		t.Error("Different timestamps should not be equal")
	}

	if !content2.Timestamp.After(content.Timestamp) {
		t.Error("Later timestamp should be after earlier timestamp")
	}
}

func TestStackFrameOperations(t *testing.T) {
	// Test stack frame operations
	frames := []StackFrame{
		{Function: "func1", File: "file1.js", Line: 10},
		{Function: "func2", File: "file2.js", Line: 20},
		{Function: "func3", File: "file3.js", Line: 30},
	}

	// Test appending
	frames = append(frames, StackFrame{Function: "func4", File: "file4.js", Line: 40})
	if len(frames) != 4 {
		t.Errorf("Expected 4 frames after append, got %d", len(frames))
	}

	// Test modification
	frames[0].Column = 5
	if frames[0].Column != 5 {
		t.Errorf("Expected column 5, got %d", frames[0].Column)
	}

	// Test slice operations
	firstTwo := frames[:2]
	if len(firstTwo) != 2 {
		t.Errorf("Expected 2 frames in slice, got %d", len(firstTwo))
	}
	if firstTwo[1].Function != "func2" {
		t.Errorf("Expected second frame function 'func2', got %q", firstTwo[1].Function)
	}
}
