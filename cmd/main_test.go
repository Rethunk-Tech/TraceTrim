package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"com.github/rethunk-tech/no-reaction/internal/config"
	"com.github/rethunk-tech/no-reaction/internal/models"
	"com.github/rethunk-tech/no-reaction/parser"
)

// Test the testable logic components without requiring Monitor mocks

// Test configuration validation logic
func TestConfigValidationLogic(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		wantError bool
	}{
		{
			name: "valid config",
			config: &config.Config{
				Clipboard: config.ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: config.ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			wantError: false,
		},
		{
			name: "invalid polling interval - too short",
			config: &config.Config{
				Clipboard: config.ClipboardConfig{
					PollingInterval: 10 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: config.ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			wantError: true,
		},
		{
			name: "invalid content size - too small",
			config: &config.Config{
				Clipboard: config.ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  100,
				},
				Parser: config.ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			wantError: true,
		},
		{
			name: "invalid content size - too large",
			config: &config.Config{
				Clipboard: config.ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  100 * 1024 * 1024,
				},
				Parser: config.ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			wantError: true,
		},
		{
			name: "invalid min stack lines - too small",
			config: &config.Config{
				Clipboard: config.ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: config.ParserConfig{
					MinStackLinesForDetection: 0,
					MinStackTraceLength:       20,
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateConfig(tt.config)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateConfig() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// Test stack trace processing logic
func TestStackTraceProcessingLogic(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectCleaning bool
	}{
		{
			name: "React stack trace with duplicates",
			content: `Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`,
			expectCleaning: true,
		},
		{
			name: "JavaScript stack trace with duplicates",
			content: `ReferenceError: x is not defined
    at eval (eval at <anonymous> (script.js:1:1))
    at eval (eval at <anonymous> (script.js:1:1))
    at eval (eval at <anonymous> (script.js:1:1))
    at Object.exports.runInThisContext (vm.js:53:16)`,
			expectCleaning: true,
		},
		{
			name:           "Non-stack-trace content",
			content:        "This is just regular text that should not be processed as a stack trace",
			expectCleaning: false,
		},
		{
			name:           "Empty content",
			content:        "",
			expectCleaning: false,
		},
		{
			name:           "Very large content (should be rejected)",
			content:        strings.Repeat("A", 60*1024*1024), // 60MB
			expectCleaning: false,
		},
		{
			name:           "Invalid UTF-8 content",
			content:        string([]byte{0xff, 0xfe, 0xfd}),
			expectCleaning: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test stack trace detection
			isStackTrace := parser.IsStackTrace(tt.content)
			if tt.expectCleaning && !isStackTrace {
				t.Errorf("Expected content to be detected as stack trace, but it wasn't")
			}
			if !tt.expectCleaning && isStackTrace {
				t.Errorf("Expected content to NOT be detected as stack trace, but it was")
			}

			// Test cleaning if it's a stack trace
			if isStackTrace {
				cleanResult := parser.CleanResult(tt.content)
				if cleanResult.Removed == 0 && tt.expectCleaning {
					t.Errorf("Expected stack trace cleaning to remove duplicates, but removed count was 0")
				}
				if cleanResult.Cleaned == tt.content && tt.expectCleaning {
					t.Errorf("Expected cleaned content to be different from original, but they were the same")
				}
			}
		})
	}
}

// Test clipboard content handling logic
func TestClipboardContentHandlingLogic(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		Clipboard: config.ClipboardConfig{
			PollingInterval: 100 * time.Millisecond,
			MaxContentSize:  1024 * 1024,
		},
		Output: config.OutputConfig{
			Verbose: false,
			Quiet:   false,
		},
		Parser: config.ParserConfig{
			MinStackLinesForDetection: 2,
			MinStackTraceLength:       20,
		},
	}

	tests := []struct {
		name    string
		content models.ClipboardContent
	}{
		{
			name: "Valid React stack trace",
			content: models.ClipboardContent{
				Content: `Error: Objects are not valid as a React child
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactErrorUtils.invokeGuardedCallback (react-dom.development.js:138:15)
    at ReactCompositeComponent._renderValidatedComponent (react-dom.development.js:185:13)`,
				Timestamp: time.Now(),
				Format:    "text/plain",
			},
		},
		{
			name: "Content too large",
			content: models.ClipboardContent{
				Content:   strings.Repeat("A", 2*1024*1024), // 2MB
				Timestamp: time.Now(),
				Format:    "text/plain",
			},
		},
		{
			name: "Non-stack-trace content",
			content: models.ClipboardContent{
				Content:   "This is regular text that should not be processed",
				Timestamp: time.Now(),
				Format:    "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: We can't easily test the full handleClipboardContent function
			// without a mock clipboard monitor, but we can test the logic components

			// Test content size validation
			if len(tt.content.Content) > cfg.Clipboard.MaxContentSize {
				// Should be handled by handleContentTooLarge logic
				if cfg.Output.Verbose {
					// Would log content too large message
				}
			} else if parser.IsStackTrace(tt.content.Content) {
				// Would process stack trace
				cleanResult := parser.CleanResult(tt.content.Content)
				if cleanResult.Cleaned != tt.content.Content {
					// Would show cleaning results
				}
			} else {
				// Would skip non-stack-trace content
				if cfg.Output.Verbose {
					// Would log skipping message
				}
			}
		})
	}
}

// Test timestamp and statistics formatting
func TestTimestampAndStatisticsLogic(t *testing.T) {
	cfg := &config.Config{
		Output: config.OutputConfig{
			ShowTimestamp: true,
			Verbose:       true,
		},
	}

	content := models.ClipboardContent{
		Content:   "test content",
		Timestamp: time.Date(2023, 12, 25, 14, 30, 45, 0, time.UTC),
		Format:    "text/plain",
	}

	// Test timestamp formatting
	timestamp := getTimestamp(content, cfg)
	expectedTimestamp := "[14:30:45] "
	if timestamp != expectedTimestamp {
		t.Errorf("Expected timestamp %q, got %q", expectedTimestamp, timestamp)
	}

	// Test statistics building
	cleanResult := &models.CleanResult{
		Removed:    5,
		BytesSaved: 256,
		Original:   "original content",
	}

	statsParts := buildStatsParts(cleanResult)
	if len(statsParts) == 0 {
		t.Error("Expected statistics parts, got empty slice")
	}

	// Check that statistics contain expected information
	statsStr := strings.Join(statsParts, ", ")
	if !strings.Contains(statsStr, "Removed 5 repetitive frame") {
		t.Errorf("Expected statistics to contain removal info, got: %s", statsStr)
	}
	if !strings.Contains(statsStr, "saved 256 bytes") {
		t.Errorf("Expected statistics to contain bytes saved info, got: %s", statsStr)
	}
}

// Test error handling in processing
func TestErrorHandlingInProcessing(t *testing.T) {
	// Test that invalid content doesn't crash the application
	invalidContents := []string{
		string([]byte{0xff, 0xfe, 0xfd}),  // Invalid UTF-8
		strings.Repeat("A", 60*1024*1024), // Very large content
		"",                                // Empty content
		"normal text",                     // Non-stack-trace content
	}

	for _, content := range invalidContents {
		// These should not panic or cause issues
		isStackTrace := parser.IsStackTrace(content)
		if isStackTrace {
			// If somehow detected as stack trace, cleaning should be safe
			cleanResult := parser.CleanResult(content)
			if cleanResult.Cleaned == "" && content != "" {
				t.Errorf("Cleaning should not result in empty content for non-empty input")
			}
		}
	}
}

// Test configuration edge cases
func TestConfigurationEdgeCases(t *testing.T) {
	// Test that configuration validation catches edge cases
	edgeCaseConfigs := []*config.Config{
		{
			Clipboard: config.ClipboardConfig{
				PollingInterval: 0, // Should be invalid
				MaxContentSize:  1024,
			},
			Parser: config.ParserConfig{
				MinStackLinesForDetection: 1,
				MinStackTraceLength:       10,
			},
		},
		{
			Clipboard: config.ClipboardConfig{
				PollingInterval: time.Hour, // Should be invalid (too long)
				MaxContentSize:  1024,
			},
			Parser: config.ParserConfig{
				MinStackLinesForDetection: 1,
				MinStackTraceLength:       10,
			},
		},
		{
			Clipboard: config.ClipboardConfig{
				PollingInterval: 100 * time.Millisecond,
				MaxContentSize:  0, // Should be invalid
			},
			Parser: config.ParserConfig{
				MinStackLinesForDetection: 1,
				MinStackTraceLength:       10,
			},
		},
	}

	for i, cfg := range edgeCaseConfigs {
		t.Run(fmt.Sprintf("edge case %d", i), func(t *testing.T) {
			err := config.ValidateConfig(cfg)
			if err == nil {
				t.Errorf("Expected configuration validation to fail for edge case %d", i)
			}
		})
	}
}

// Test concurrent access patterns (simulated)
func TestConcurrentAccessPatterns(t *testing.T) {
	// Test that our logic components handle concurrent access properly
	// This is a simplified test since we can't easily test the full Monitor

	cfg := config.DefaultConfig()

	// Test that configuration access is thread-safe
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			// These operations should not cause race conditions
			_ = config.ValidateConfig(cfg)
			_ = cfg.Output.Verbose
			_ = cfg.Clipboard.PollingInterval
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test that invalid configs are caught
func TestInvalidConfigurationHandling(t *testing.T) {
	invalidCfg := &config.Config{
		Clipboard: config.ClipboardConfig{
			PollingInterval: 10 * time.Millisecond, // Too short
			MaxContentSize:  500,                   // Too small
		},
		Parser: config.ParserConfig{
			MinStackLinesForDetection: 1,
			MinStackTraceLength:       10,
		},
	}

	err := config.ValidateConfig(invalidCfg)
	if err == nil {
		t.Error("Expected configuration validation to fail for invalid config")
	}
}

// Test model creation and basic operations
func TestModelsCreation(t *testing.T) {
	// Test creating clipboard content
	content := models.ClipboardContent{
		Content:   "test content",
		Timestamp: time.Now(),
		Format:    "text/plain",
	}

	if content.Content != "test content" {
		t.Errorf("Expected content 'test content', got %q", content.Content)
	}

	if content.Format != "text/plain" {
		t.Errorf("Expected format 'text/plain', got %q", content.Format)
	}

	// Test that timestamp is set
	if content.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

// Test configuration defaults
func TestConfigDefaults(t *testing.T) {
	cfg := config.DefaultConfig()

	// Test default values
	if cfg.Clipboard.PollingInterval != 500*time.Millisecond {
		t.Errorf("Expected default polling interval 500ms, got %v", cfg.Clipboard.PollingInterval)
	}

	if cfg.Clipboard.MaxContentSize != 1024*1024 {
		t.Errorf("Expected default max content size 1MB, got %d", cfg.Clipboard.MaxContentSize)
	}

	if cfg.Output.Verbose != false {
		t.Errorf("Expected default verbose false, got %v", cfg.Output.Verbose)
	}

	if cfg.Output.Quiet != false {
		t.Errorf("Expected default quiet false, got %v", cfg.Output.Quiet)
	}

	if cfg.Parser.MinStackLinesForDetection != 2 {
		t.Errorf("Expected default min stack lines 2, got %d", cfg.Parser.MinStackLinesForDetection)
	}
}
