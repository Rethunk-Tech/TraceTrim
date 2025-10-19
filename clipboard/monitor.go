package clipboard

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"com.github/rethunk-tech/tracetrim/internal/models"
	"golang.design/x/clipboard"
)

const (
	// Clipboard polling interval
	clipboardPollInterval = 500 * time.Millisecond
)

// Monitor handles clipboard monitoring across platforms
type Monitor struct {
	platform    Platform
	stopChan    chan struct{}
	lastContent string
	mutex       sync.RWMutex // Protects lastContent
}

// Platform interface abstracts platform-specific clipboard operations
type Platform interface {
	GetContent() (string, error)
	SetContent(content string) error
	GetName() string
}

// standardPlatform implements Platform interface using golang.design/x/clipboard
type standardPlatform struct{}

// GetName returns the platform name
func (s *standardPlatform) GetName() string {
	return "Standard"
}

// GetContent retrieves text content from clipboard using standard library
func (s *standardPlatform) GetContent() (string, error) {
	// Initialize clipboard if needed
	err := clipboard.Init()
	if err != nil {
		return "", fmt.Errorf("failed to initialize clipboard: %w", err)
	}

	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return "", fmt.Errorf("no text data available in clipboard")
	}

	return string(data), nil
}

// SetContent sets text content to clipboard using standard library
func (s *standardPlatform) SetContent(content string) error {
	// Initialize clipboard if needed
	err := clipboard.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize clipboard: %w", err)
	}

	clipboard.Write(clipboard.FmtText, []byte(content))
	return nil
}

// NewMonitor creates a new clipboard monitor for the current platform
func NewMonitor() (*Monitor, error) {
	platform, err := getPlatform()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize clipboard platform: %w", err)
	}

	return &Monitor{
		platform: platform,
		stopChan: make(chan struct{}),
	}, nil
}

// StartMonitoring begins monitoring the clipboard for changes with default interval
func (m *Monitor) StartMonitoring(ctx context.Context, callback func(models.ClipboardContent, *Monitor)) error {
	return m.StartMonitoringWithInterval(ctx, clipboardPollInterval, callback)
}

// StartMonitoringWithInterval begins monitoring the clipboard for changes with custom interval
func (m *Monitor) StartMonitoringWithInterval(ctx context.Context, interval time.Duration, callback func(models.ClipboardContent, *Monitor)) error {
	log.Printf("Starting clipboard monitoring on %s with %v interval", m.platform.GetName(), interval)

	// Get initial content
	initialContent, err := m.platform.GetContent()
	if err != nil {
		return fmt.Errorf("failed to get initial clipboard content: %w", err)
	}
	m.mutex.Lock()
	m.lastContent = initialContent
	m.mutex.Unlock()

	// Start monitoring loop
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping clipboard monitoring")
			return nil
		case <-m.stopChan:
			log.Println("Stopping clipboard monitoring")
			return nil
		case <-ticker.C:
			content, err := m.platform.GetContent()
			if err != nil {
				log.Printf("Error getting clipboard content: %v", err)
				continue
			}

			// Check if content has changed (with proper locking)
			m.mutex.Lock()
			contentChanged := content != m.lastContent && content != ""
			if contentChanged {
				m.lastContent = content
			}
			m.mutex.Unlock()

			if contentChanged {
				clipboardContent := models.ClipboardContent{
					Content:   content,
					Timestamp: time.Now(),
					Format:    "text/plain",
				}
				callback(clipboardContent, m)
			}
		}
	}
}

// Stop stops the clipboard monitoring
func (m *Monitor) Stop() {
	select {
	case m.stopChan <- struct{}{}:
	default:
	}
}

// GetCurrentContent returns the current clipboard content
func (m *Monitor) GetCurrentContent() (string, error) {
	return m.platform.GetContent()
}

// SetContent sets the clipboard content
func (m *Monitor) SetContent(content string) error {
	return m.platform.SetContent(content)
}
