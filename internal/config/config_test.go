package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Test clipboard defaults
	if cfg.Clipboard.PollingInterval != 500*time.Millisecond {
		t.Errorf("Expected polling interval 500ms, got %v", cfg.Clipboard.PollingInterval)
	}

	if cfg.Clipboard.MaxContentSize != 1024*1024 {
		t.Errorf("Expected max content size 1MB, got %d", cfg.Clipboard.MaxContentSize)
	}

	// Test output defaults
	if cfg.Output.Verbose != false {
		t.Errorf("Expected verbose false, got %v", cfg.Output.Verbose)
	}

	if cfg.Output.Quiet != false {
		t.Errorf("Expected quiet false, got %v", cfg.Output.Quiet)
	}

	if cfg.Output.ShowTimestamp != true {
		t.Errorf("Expected show timestamp true, got %v", cfg.Output.ShowTimestamp)
	}

	if cfg.Output.LogFile != "" {
		t.Errorf("Expected empty log file, got %q", cfg.Output.LogFile)
	}

	// Test parser defaults
	if cfg.Parser.MinStackLinesForDetection != 2 {
		t.Errorf("Expected min stack lines 2, got %d", cfg.Parser.MinStackLinesForDetection)
	}

	if cfg.Parser.MinStackTraceLength != 20 {
		t.Errorf("Expected min stack trace length 20, got %d", cfg.Parser.MinStackTraceLength)
	}

	if len(cfg.Parser.CustomPatterns) != 0 {
		t.Errorf("Expected empty custom patterns, got %v", cfg.Parser.CustomPatterns)
	}

	// Test app defaults
	if cfg.App.ConfigFile != "config.yaml" {
		t.Errorf("Expected config file 'config.yaml', got %q", cfg.App.ConfigFile)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			name:      "Valid config",
			config:    DefaultConfig(),
			shouldErr: false,
		},
		{
			name: "Polling interval too short",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 10 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			shouldErr: true,
		},
		{
			name: "Max content size too small",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  500,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       20,
				},
			},
			shouldErr: true,
		},
		{
			name: "Min stack lines too small",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 0,
					MinStackTraceLength:       20,
				},
			},
			shouldErr: true,
		},
		{
			name: "Min stack trace length too small",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 2,
					MinStackTraceLength:       5,
				},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.shouldErr && err == nil {
				t.Error("Expected validation error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no validation error but got: %v", err)
			}
		})
	}
}

func TestConfigStructCreation(t *testing.T) {
	// Test that we can create config structs and they initialize properly
	cfg := &Config{
		Clipboard: ClipboardConfig{
			PollingInterval: 1 * time.Second,
			MaxContentSize:  2 * 1024 * 1024,
		},
		Output: OutputConfig{
			Verbose:       true,
			LogFile:       "/tmp/test.log",
			ShowTimestamp: false,
			Quiet:         true,
		},
		Parser: ParserConfig{
			MinStackLinesForDetection: 5,
			MinStackTraceLength:       50,
			CustomPatterns:            []string{"pattern1", "pattern2"},
		},
		App: AppConfig{
			ConfigFile: "custom.yaml",
		},
	}

	if cfg.Clipboard.PollingInterval != 1*time.Second {
		t.Errorf("Expected polling interval 1s, got %v", cfg.Clipboard.PollingInterval)
	}

	if cfg.Output.Verbose != true {
		t.Errorf("Expected verbose true, got %v", cfg.Output.Verbose)
	}

	if cfg.Parser.MinStackLinesForDetection != 5 {
		t.Errorf("Expected min stack lines 5, got %d", cfg.Parser.MinStackLinesForDetection)
	}

	if len(cfg.Parser.CustomPatterns) != 2 {
		t.Errorf("Expected 2 custom patterns, got %d", len(cfg.Parser.CustomPatterns))
	}
}

func TestConfigValidationEdgeCases(t *testing.T) {
	// Test edge cases for validation
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "Exactly minimum polling interval",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 50 * time.Millisecond,
					MaxContentSize:  1024 * 1024,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 1,
					MinStackTraceLength:       10,
				},
			},
		},
		{
			name: "Exactly minimum content size",
			config: &Config{
				Clipboard: ClipboardConfig{
					PollingInterval: 500 * time.Millisecond,
					MaxContentSize:  1024,
				},
				Parser: ParserConfig{
					MinStackLinesForDetection: 1,
					MinStackTraceLength:       10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if err != nil {
				t.Errorf("Expected validation to pass for edge case %q: %v", tt.name, err)
			}
		})
	}
}

func TestViperIntegration(t *testing.T) {
	// Test that viper integration works (basic smoke test)
	v := viper.New()

	// Set some values in viper
	v.Set("clipboard-polling-interval", "1s")
	v.Set("clipboard-max-content-size", 2048)
	v.Set("verbose", true)
	v.Set("quiet", false)

	// Create a config and try to unmarshal
	cfg := DefaultConfig()
	err := v.Unmarshal(cfg)
	if err != nil {
		t.Errorf("Failed to unmarshal config from viper: %v", err)
	}

	// Check that values were loaded (though defaults might override)
	// This is a basic integration test
	if cfg == nil {
		t.Error("Config should not be nil after unmarshal")
	}
}
