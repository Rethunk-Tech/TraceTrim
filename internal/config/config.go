package config

import (
	"fmt"
	"regexp"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// Default clipboard polling interval
	DefaultPollingInterval = 500 * time.Millisecond

	// Default maximum clipboard content size (1MB)
	DefaultMaxContentSize = 1024 * 1024

	// Minimum clipboard content size (1KB)
	MinContentSize = 1024

	// Maximum clipboard content size (50MB)
	MaxContentSize = 50 * 1024 * 1024

	// Maximum clipboard polling interval (10 seconds)
	MaxPollingInterval = 10 * time.Second

	// Minimum parser stack lines for detection
	MinStackLines = 1

	// Maximum parser stack lines for detection
	MaxStackLines = 100

	// Minimum parser stack trace length
	MinStackTraceLength = 10

	// Maximum parser stack trace length
	MaxStackTraceLength = 10000

	// Default parser minimum stack lines for detection
	DefaultMinStackLines = 2

	// Default parser minimum stack trace length
	DefaultMinStackTraceLength = 20
)

// Config holds all configuration for the application
type Config struct {
	// Application settings
	App AppConfig

	// Parser settings
	Parser ParserConfig

	// Output and logging settings
	Output OutputConfig

	// Script mode settings
	Script ScriptConfig

	// Clipboard monitoring settings
	Clipboard ClipboardConfig

	// Script mode flag (simplified)
	ScriptMode bool

	// Auto-detect script mode based on environment
	AutoDetectScriptMode bool
}

// ClipboardConfig contains clipboard monitoring configuration
type ClipboardConfig struct {
	// PollingInterval is how often to check for clipboard changes
	PollingInterval time.Duration `mapstructure:"clipboard-polling-interval"`

	// MaxContentSize is the maximum clipboard content size to process (in bytes)
	MaxContentSize int `mapstructure:"clipboard-max-content-size"`
}

// OutputConfig contains output and logging configuration
type OutputConfig struct {
	// LogFile is the path to log file (empty for stdout)
	LogFile string

	// Verbose enables detailed logging
	Verbose bool

	// ShowTimestamp controls whether to show timestamps in output
	ShowTimestamp bool

	// Quiet suppresses non-essential output
	Quiet bool
}

// ParserConfig contains parser-specific configuration
type ParserConfig struct {
	// CustomPatterns allows adding custom regex patterns for stack trace detection
	CustomPatterns []string

	// MinStackLinesForDetection minimum lines to consider content a stack trace
	MinStackLinesForDetection int

	// MinStackTraceLength minimum content length to consider for stack trace detection
	MinStackTraceLength int
}

// ScriptConfig contains script mode configuration
type ScriptConfig struct {
	// OutputFormat controls the output format in script mode
	OutputFormat string

	// Enabled determines if script mode is active
	Enabled bool

	// ShowStatistics controls whether to show cleaning statistics in script mode
	ShowStatistics bool

	// ExitCodeOnError controls whether to exit with error code when no stack trace is found
	ExitCodeOnError bool
}

// AppConfig contains general application configuration
type AppConfig struct {
	// ConfigFile path to configuration file
	ConfigFile string
}

// DefaultConfig returns default configuration values
func DefaultConfig() *Config {
	return &Config{
		Clipboard: ClipboardConfig{
			PollingInterval: DefaultPollingInterval,
			MaxContentSize:  DefaultMaxContentSize, // 1MB
		},
		Output: OutputConfig{
			Verbose:       false,
			LogFile:       "",
			ShowTimestamp: true,
			Quiet:         false,
		},
		Parser: ParserConfig{
			MinStackLinesForDetection: DefaultMinStackLines,
			MinStackTraceLength:       DefaultMinStackTraceLength,
			CustomPatterns:            []string{},
		},
		Script: ScriptConfig{
			Enabled:         false,
			OutputFormat:    "cleaned", // "cleaned", "json", "stats"
			ShowStatistics:  true,
			ExitCodeOnError: false,
		},
		App: AppConfig{
			ConfigFile: "config.yaml",
		},
		ScriptMode:           false,
		AutoDetectScriptMode: true, // Enable auto-detection by default
	}
}

// LoadConfig loads configuration from file and command line flags
func LoadConfig() (*Config, error) {
	// Use the global viper instance
	v := viper.GetViper()

	// Set configuration file properties
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Enable reading from environment variables
	v.SetEnvPrefix("NO_REACTION")
	v.AutomaticEnv()

	// Try to read config file first (ignore error if file doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		// Log that we're using defaults if not in quiet mode
		// We need to check if quiet flag is set, but we don't have config yet
		// So we'll check the flag directly
		if !v.GetBool("quiet") {
			fmt.Printf("No config file found, using defaults\n")
		}
	}

	// Start with default config
	config := DefaultConfig()

	// Try to unmarshal first to get config file values
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Override with flag values if they were explicitly set
	// This ensures flags take precedence over config file
	if v.IsSet("clipboard-polling-interval") {
		config.Clipboard.PollingInterval = v.GetDuration("clipboard-polling-interval")
	}
	if v.IsSet("clipboard-max-content-size") {
		config.Clipboard.MaxContentSize = v.GetInt("clipboard-max-content-size")
	}
	if v.IsSet("verbose") {
		config.Output.Verbose = v.GetBool("verbose")
	}
	if v.IsSet("log-file") {
		config.Output.LogFile = v.GetString("log-file")
	}
	if v.IsSet("quiet") {
		config.Output.Quiet = v.GetBool("quiet")
	}
	if v.IsSet("show-timestamp") {
		config.Output.ShowTimestamp = v.GetBool("show-timestamp")
	}
	if v.IsSet("parser-min-stack-lines") {
		config.Parser.MinStackLinesForDetection = v.GetInt("parser-min-stack-lines")
	}
	if v.IsSet("parser-min-stack-trace-length") {
		config.Parser.MinStackTraceLength = v.GetInt("parser-min-stack-trace-length")
	}
	if v.IsSet("parser-custom-patterns") {
		config.Parser.CustomPatterns = v.GetStringSlice("parser-custom-patterns")
	}
	if v.IsSet("script-mode") {
		config.ScriptMode = v.GetBool("script-mode")
	}
	if v.IsSet("auto-detect-script-mode") {
		config.AutoDetectScriptMode = v.GetBool("auto-detect-script-mode")
	}

	return config, nil
}

// BindFlags binds command line flags to viper
func BindFlags() error {
	pflag.String("config", "config.yaml", "Configuration file path")
	pflag.Duration("clipboard-polling-interval", DefaultPollingInterval, "Clipboard polling interval")
	pflag.Int("clipboard-max-content-size", DefaultMaxContentSize, "Maximum clipboard content size in bytes")
	pflag.Bool("verbose", false, "Enable verbose output")
	pflag.String("log-file", "", "Log file path (empty for stdout)")
	pflag.Bool("quiet", false, "Suppress non-essential output")
	pflag.Bool("show-timestamp", true, "Show timestamps in output")
	pflag.Int("parser-min-stack-lines", DefaultMinStackLines, "Minimum stack lines for detection")
	pflag.Int("parser-min-stack-trace-length", DefaultMinStackTraceLength, "Minimum stack trace length")
	pflag.StringSlice("parser-custom-patterns", []string{}, "Custom regex patterns for stack trace detection")
	pflag.Bool("script-mode", false, "Enable script mode (read from STDIN, write to STDOUT, then exit)")
	pflag.Bool("auto-detect-script-mode", true, "Auto-detect script mode based on non-interactive environment")

	// Bind flags to the global viper instance FIRST
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}

	// Parse flags AFTER binding to viper
	pflag.Parse()

	return nil
}

// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	// Validate clipboard configuration
	if config.Clipboard.PollingInterval < 50*time.Millisecond {
		return fmt.Errorf("clipboard polling interval must be at least 50ms")
	}
	if config.Clipboard.PollingInterval > MaxPollingInterval {
		return fmt.Errorf("clipboard polling interval must be at most 10 seconds")
	}

	if config.Clipboard.MaxContentSize < MinContentSize {
		return fmt.Errorf("clipboard max content size must be at least 1KB")
	}
	if config.Clipboard.MaxContentSize > MaxContentSize {
		return fmt.Errorf("clipboard max content size must be at most 50MB")
	}

	// Validate parser configuration
	if config.Parser.MinStackLinesForDetection < MinStackLines {
		return fmt.Errorf("parser min stack lines must be at least 1")
	}
	if config.Parser.MinStackLinesForDetection > MaxStackLines {
		return fmt.Errorf("parser min stack lines must be at most 100")
	}

	if config.Parser.MinStackTraceLength < MinStackTraceLength {
		return fmt.Errorf("parser min stack trace length must be at least 10")
	}
	if config.Parser.MinStackTraceLength > MaxStackTraceLength {
		return fmt.Errorf("parser min stack trace length must be at most 10000")
	}

	// Validate custom patterns if provided
	for i, pattern := range config.Parser.CustomPatterns {
		if pattern == "" {
			return fmt.Errorf("custom pattern at index %d cannot be empty", i)
		}
		// Validate that pattern is a valid regex by attempting to compile it
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("custom pattern at index %d is not a valid regex: %w", i, err)
		}
	}

	return nil
}
