package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"com.github/rethunk-tech/no-reaction/clipboard"
	"com.github/rethunk-tech/no-reaction/internal/config"
	"com.github/rethunk-tech/no-reaction/internal/models"
	"com.github/rethunk-tech/no-reaction/parser"
)

func main() {
	// Bind command line flags to viper
	if err := config.BindFlags(); err != nil {
		log.Fatalf("Failed to bind flags: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Set up logging based on configuration
	if cfg.Output.LogFile != "" {
		// TODO: Implement file logging if needed
	}

	// Print startup information based on verbosity
	if !cfg.Output.Quiet {
		fmt.Println("Clipboard Stack Trace Cleaner")
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
		log.Fatalf("Failed to initialize clipboard monitor: %v", err)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring in a goroutine
	go func() {
		err := monitor.StartMonitoringWithInterval(ctx, cfg.Clipboard.PollingInterval, func(content models.ClipboardContent, m *clipboard.Monitor) {
			handleClipboardContent(content, m, cfg)
		})
		if err != nil {
			log.Printf("Monitoring stopped with error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down...")
	monitor.Stop()
}

// handleClipboardContent processes clipboard content when it changes
func handleClipboardContent(content models.ClipboardContent, monitor *clipboard.Monitor, cfg *config.Config) {
	// Check content size limit
	if len(content.Content) > cfg.Clipboard.MaxContentSize {
		if cfg.Output.Verbose {
			log.Printf("Content too large (%d bytes), skipping", len(content.Content))
		}
		return
	}

	// Check if this looks like a stack trace
	if parser.IsStackTrace(content.Content) {
		timestamp := ""
		if cfg.Output.ShowTimestamp {
			timestamp = fmt.Sprintf("[%s] ", content.Timestamp.Format("15:04:05"))
		}

		if cfg.Output.Verbose {
			fmt.Printf("%sDetected stack trace, cleaning...\n", timestamp)
		}

		// Clean the stack trace and get detailed results
		cleanResult := parser.CleanResult(content.Content)

		// Check if content actually changed
		if cleanResult.Cleaned == content.Content {
			if cfg.Output.Verbose {
				fmt.Printf("%sNo changes needed - content is already clean\n", timestamp)
			}
			return
		}

		// Update clipboard with cleaned content using the existing monitor
		err := monitor.SetContent(cleanResult.Cleaned)
		if err != nil {
			log.Printf("Failed to update clipboard: %v", err)
			return
		}

		if !cfg.Output.Quiet {
			fmt.Printf("%s✓ Stack trace cleaned and clipboard updated\n", timestamp)
		}

		// Show detailed statistics in verbose mode
		if cfg.Output.Verbose {
			fmt.Printf("%sStatistics:\n", timestamp)
			fmt.Printf("%s  • Size: %d bytes → %d bytes", timestamp, len(cleanResult.Original), len(cleanResult.Cleaned))

			if len(cleanResult.Original) > len(cleanResult.Cleaned) {
				bytesSaved := len(cleanResult.Original) - len(cleanResult.Cleaned)
				fmt.Printf(" (saved %d bytes", bytesSaved)
				if len(cleanResult.Original) > 0 {
					percentage := float64(bytesSaved) / float64(len(cleanResult.Original)) * 100
					fmt.Printf(", %.1f%%", percentage)
				}
				fmt.Printf(")")
			}
			fmt.Printf("\n")

			if cleanResult.Removed > 0 {
				fmt.Printf("%s  • Removed %d repetitive stack frame(s)\n", timestamp, cleanResult.Removed)
			}
		}
	}
}
