package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"com.github/rethunk-tech/no-reaction/clipboard"
	"com.github/rethunk-tech/no-reaction/internal/models"
	"com.github/rethunk-tech/no-reaction/parser"
)

func main() {
	fmt.Println("Clipboard Stack Trace Cleaner")
	fmt.Println("Monitoring clipboard for JavaScript/React stack traces...")
	fmt.Println("Press Ctrl+C to exit")

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
		err := monitor.StartMonitoring(ctx, func(content models.ClipboardContent) {
			handleClipboardContent(content)
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
func handleClipboardContent(content models.ClipboardContent) {
	// Check if this looks like a stack trace
	if parser.IsStackTrace(content.Content) {
		fmt.Printf("\n[%s] Detected stack trace, cleaning...\n", content.Timestamp.Format("15:04:05"))

		// Clean the stack trace
		cleaned := parser.CleanStackTrace(content.Content)

		// Update clipboard with cleaned content
		monitor, err := clipboard.NewMonitor()
		if err != nil {
			log.Printf("Failed to create monitor for updating clipboard: %v", err)
			return
		}

		err = monitor.SetContent(cleaned)
		if err != nil {
			log.Printf("Failed to update clipboard: %v", err)
			return
		}

		fmt.Printf("✓ Stack trace cleaned and clipboard updated\n")

		// Show a brief preview of what was changed
		originalLines := strings.Split(content.Content, "\n")
		cleanedLines := strings.Split(cleaned, "\n")

		if len(cleanedLines) < len(originalLines) {
			removed := len(originalLines) - len(cleanedLines)
			fmt.Printf("  Removed %d repetitive lines\n", removed)
		}
	}
}
