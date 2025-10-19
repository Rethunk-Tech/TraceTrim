package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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
		err := monitor.StartMonitoring(ctx, func(content models.ClipboardContent, m *clipboard.Monitor) {
			handleClipboardContent(content, m)
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
func handleClipboardContent(content models.ClipboardContent, monitor *clipboard.Monitor) {
	// Check if this looks like a stack trace
	if parser.IsStackTrace(content.Content) {
		fmt.Printf("\n[%s] Detected stack trace, cleaning...\n", content.Timestamp.Format("15:04:05"))

		// Clean the stack trace and get detailed results
		cleanResult := parser.CleanResult(content.Content)

		// Check if content actually changed
		if cleanResult.Cleaned == content.Content {
			fmt.Printf("  No changes needed - content is already clean\n")
			return
		}

		// Update clipboard with cleaned content using the existing monitor
		err := monitor.SetContent(cleanResult.Cleaned)
		if err != nil {
			log.Printf("Failed to update clipboard: %v", err)
			return
		}

		fmt.Printf("✓ Stack trace cleaned and clipboard updated\n")

		// Show accurate count of what was removed
		if cleanResult.Removed > 0 {
			fmt.Printf("  Removed %d repetitive stack frame(s)\n", cleanResult.Removed)
		}
	}
}
