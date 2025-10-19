package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"com.github/rethunk-tech/tracetrim/clipboard"
	"com.github/rethunk-tech/tracetrim/internal/config"
	"com.github/rethunk-tech/tracetrim/internal/models"
	"com.github/rethunk-tech/tracetrim/parser"
)

// version is set during build time via ldflags
var version = "dev"

// Constants for stack trace types
const (
	stackTraceTypeReact      = "React"
	stackTraceTypeJavaScript = "JavaScript"
)

func main() {
	// Bind command line flags to viper
	if err := config.BindFlags(); err != nil {
		log.Fatalf("Failed to bind flags: %v", err)
	}

	// Load configuration
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatalf("Failed to load configuration: %v", configErr)
	}

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Set up logging based on configuration
	// Note: File logging is not implemented in this version

	// Print startup information based on verbosity
	if !cfg.Output.Quiet {
		fmt.Printf("TraceTrim v%s\n", version)
		if cfg.Output.Verbose {
			fmt.Printf("Configuration loaded from: %s\n", cfg.App.ConfigFile)
			fmt.Printf("Polling interval: %v\n", cfg.Clipboard.PollingInterval)
		}
		fmt.Println("Monitoring clipboard for JavaScript/React stack traces...")
		fmt.Println("Press Ctrl+C to exit")
	}

	// Create clipboard monitor
	monitor, err := clipboard.NewMonitor()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize clipboard monitor: %v\n", err)
		fmt.Fprintf(os.Stderr, "This may be due to:\n")
		fmt.Fprintf(os.Stderr, "  - Insufficient permissions to access clipboard\n")
		fmt.Fprintf(os.Stderr, "  - Platform-specific requirements not met\n")
		fmt.Fprintf(os.Stderr, "  - Missing system dependencies\n")
		fmt.Fprintf(os.Stderr, "\nPlease check the troubleshooting section in the README.\n")
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring in a goroutine
	go func() {
		callback := func(content models.ClipboardContent, m *clipboard.Monitor) {
			handleClipboardContent(content, m, cfg)
		}
		err := monitor.StartMonitoringWithInterval(ctx, cfg.Clipboard.PollingInterval, callback)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Clipboard monitoring stopped: %v\n", err)
			// Try to restart monitoring after a delay
			time.Sleep(5 * time.Second)
			fmt.Fprintf(os.Stderr, "Info: Attempting to restart clipboard monitoring...\n")
			go func() {
				restartErr := monitor.StartMonitoringWithInterval(ctx, cfg.Clipboard.PollingInterval, callback)
				if restartErr != nil {
					fmt.Fprintf(os.Stderr, "Error: Failed to restart monitoring: %v\n", restartErr)
				}
			}()
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down...")
	monitor.Stop()
}

// plural returns "s" if count != 1, otherwise returns empty string
func plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// getStackTraceType determines the type of stack trace for better user feedback
func getStackTraceType(errorInfo *models.ErrorInfo, originalContent string) string {
	if errorInfo == nil {
		return stackTraceTypeJavaScript
	}

	// Check for React component
	if errorInfo.Component != "" {
		return stackTraceTypeReact
	}

	// Check for React-related files in source
	if errorInfo.Source != "" {
		sourceLower := strings.ToLower(errorInfo.Source)
		if strings.Contains(sourceLower, "react") ||
			strings.Contains(sourceLower, "jsx") ||
			strings.Contains(sourceLower, "tsx") {
			return stackTraceTypeReact
		}
	}

	// Check original content for React patterns
	contentLower := strings.ToLower(originalContent)
	if strings.Contains(contentLower, "react") ||
		strings.Contains(contentLower, "component") ||
		strings.Contains(contentLower, "jsx") ||
		strings.Contains(contentLower, "tsx") {
		return stackTraceTypeReact
	}

	return stackTraceTypeJavaScript
}

// handleClipboardContent processes clipboard content when it changes
func handleClipboardContent(content models.ClipboardContent, monitor *clipboard.Monitor, cfg *config.Config) {
	// Check content size limit
	if len(content.Content) > cfg.Clipboard.MaxContentSize {
		handleContentTooLarge(content, cfg)
		return
	}

	// Check if this looks like a stack trace
	if !parser.IsStackTrace(content.Content) {
		if cfg.Output.Verbose {
			timestamp := getTimestamp(content, cfg)
			fmt.Printf("%sSkipping non-stack-trace content\n", timestamp)
		}
		return
	}

	// Process stack trace
	processStackTrace(content, monitor, cfg)
}

// handleContentTooLarge logs when content is too large to process
func handleContentTooLarge(content models.ClipboardContent, cfg *config.Config) {
	if cfg.Output.Verbose {
		log.Printf("Content too large (%d bytes), skipping", len(content.Content))
	}
}

// processStackTrace handles the main stack trace processing logic
func processStackTrace(content models.ClipboardContent, monitor *clipboard.Monitor, cfg *config.Config) {
	// Clean the stack trace and get detailed results
	cleanResult := parser.CleanResult(content.Content)

	// Check if content actually changed
	if cleanResult.Cleaned == content.Content {
		handleUnchangedContent(content, cfg)
		return
	}

	// Update clipboard with cleaned content
	if err := updateClipboard(monitor, &cleanResult); err != nil {
		timestamp := getTimestamp(content, cfg)
		fmt.Fprintf(os.Stderr, "%sError: Failed to update clipboard: %v\n", timestamp, err)
		fmt.Fprintf(os.Stderr, "%sThe cleaned content could not be written back to clipboard\n", timestamp)
		return
	}

	// Show results
	showCleaningResults(content, &cleanResult, cfg)
}

// updateClipboard updates the clipboard with cleaned content
func updateClipboard(monitor *clipboard.Monitor, cleanResult *models.CleanResult) error {
	return monitor.SetContent(cleanResult.Cleaned)
}

// handleUnchangedContent handles the case where content is already clean
func handleUnchangedContent(content models.ClipboardContent, cfg *config.Config) {
	if cfg.Output.Verbose {
		timestamp := getTimestamp(content, cfg)
		fmt.Printf("%sNo changes needed - content is already clean\n", timestamp)
	}
}

// showCleaningResults displays the results of cleaning a stack trace
func showCleaningResults(content models.ClipboardContent, cleanResult *models.CleanResult, cfg *config.Config) {
	timestamp := getTimestamp(content, cfg)

	if cfg.Output.Verbose {
		stackType := getStackTraceType(cleanResult.ErrorInfo, content.Content)
		fmt.Printf("%s🔍 Detected %s stack trace, cleaning...\n", timestamp, stackType)
	}

	if !cfg.Output.Quiet {
		showSuccessMessage(content, cleanResult, cfg)
		showCompactStatistics(timestamp, cleanResult)
	}

	if cfg.Output.Verbose {
		showVerboseStatistics(timestamp, cleanResult)
	}
}

// showSuccessMessage displays the success message with stack trace type
func showSuccessMessage(content models.ClipboardContent, cleanResult *models.CleanResult, cfg *config.Config) {
	timestamp := getTimestamp(content, cfg)
	stackType := getStackTraceType(cleanResult.ErrorInfo, content.Content)
	fmt.Printf("%s✅ %s stack trace cleaned and clipboard updated\n", timestamp, stackType)
}

// getTimestamp returns formatted timestamp if enabled in config
func getTimestamp(content models.ClipboardContent, cfg *config.Config) string {
	if cfg.Output.ShowTimestamp {
		return fmt.Sprintf("[%s] ", content.Timestamp.Format("15:04:05"))
	}
	return ""
}

// showCompactStatistics displays compact statistics for cleaned content
func showCompactStatistics(timestamp string, cleanResult *models.CleanResult) {
	if cleanResult.Removed > 0 || cleanResult.BytesSaved > 0 {
		fmt.Printf("%s   • ", timestamp)

		statsParts := buildStatsParts(cleanResult)
		fmt.Printf("%s\n", strings.Join(statsParts, ", "))
	}
}

// showVerboseStatistics displays verbose statistics for cleaned content
func showVerboseStatistics(timestamp string, cleanResult *models.CleanResult) {
	fmt.Printf("%s   • ", timestamp)

	statsParts := buildStatsParts(cleanResult)
	if len(statsParts) > 0 {
		fmt.Printf("%s\n", strings.Join(statsParts, ", "))
	} else {
		fmt.Printf("No changes needed\n")
	}
}

// buildStatsParts builds the statistics parts for display
func buildStatsParts(cleanResult *models.CleanResult) []string {
	statsParts := []string{}

	if cleanResult.Removed > 0 {
		statsParts = append(statsParts, fmt.Sprintf("Removed %d repetitive frame%s", cleanResult.Removed, plural(cleanResult.Removed)))
	}

	if cleanResult.BytesSaved > 0 {
		percentage := float64(cleanResult.BytesSaved) / float64(len(cleanResult.Original)) * 100
		statsParts = append(statsParts, fmt.Sprintf("saved %d bytes, %.1f%%", cleanResult.BytesSaved, percentage))
	}

	return statsParts
}
