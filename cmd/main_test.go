package main

import (
	"testing"
	"time"

	"com.github/rethunk-tech/no-reaction/internal/config"
	"com.github/rethunk-tech/no-reaction/internal/models"
)

// Test the testable logic components without requiring Monitor mocks

// Test configuration validation logic
func TestConfigValidationLogic(t *testing.T) {
	// Test that invalid configs are caught
	invalidCfg := &config.Config{
		Clipboard: config.ClipboardConfig{
			PollingInterval: 10 * time.Millisecond, // Too short
			MaxContentSize:  500,                   // Too small
		},
		Parser: config.ParserConfig{
			MinStackLinesForDetection: 0, // Too small
			MinStackTraceLength:       5, // Too small
		},
	}

	err := config.ValidateConfig(invalidCfg)
	if err == nil {
		t.Error("Expected validation to fail for invalid config")
	}

	// Test that valid configs pass
	validCfg := config.DefaultConfig()
	err = config.ValidateConfig(validCfg)
	if err != nil {
		t.Errorf("Expected validation to pass for valid config: %v", err)
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
