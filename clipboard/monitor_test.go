package clipboard

import (
	"context"
	"testing"
	"time"

	"com.github/rethunk-tech/no-reaction/internal/models"
)

// mockPlatform implements Platform interface for testing
type mockPlatform struct {
	content       string
	setContent    string
	name          string
	callCount     int
	changeContent bool // Whether to change content on each call
}

func (m *mockPlatform) GetContent() (string, error) {
	m.callCount++
	// Change content after 2 calls to simulate a clipboard change
	// Call 1: Initial content (before monitoring starts)
	// Call 2: Initial content (first monitoring check - no callback)
	// Call 3+: Changed content (second monitoring check - callback triggered)
	if m.changeContent && m.callCount >= 3 {
		return "changed content", nil
	}
	return m.content, nil
}

func (m *mockPlatform) SetContent(content string) error {
	m.setContent = content
	return nil
}

func (m *mockPlatform) GetName() string {
	return m.name
}

func TestNewMonitor(t *testing.T) {
	monitor, err := NewMonitor()
	if err != nil {
		t.Fatalf("NewMonitor() failed: %v", err)
	}
	if monitor == nil {
		t.Fatal("NewMonitor() returned nil monitor")
	}
	if monitor.platform == nil {
		t.Fatal("NewMonitor() returned monitor with nil platform")
	}
	if monitor.stopChan == nil {
		t.Fatal("NewMonitor() returned monitor with nil stopChan")
	}
}

func TestMonitor_GetCurrentContent(t *testing.T) {
	// Use a mock platform for testing
	mockPlatform := &mockPlatform{
		content: "test content",
		name:    "test",
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	content, err := monitor.GetCurrentContent()
	if err != nil {
		t.Fatalf("GetCurrentContent() failed: %v", err)
	}
	if content != "test content" {
		t.Errorf("GetCurrentContent() = %q, want %q", content, "test content")
	}
}

func TestMonitor_SetContent(t *testing.T) {
	// Use a mock platform for testing SetContent
	mockPlatform := &mockPlatform{
		content: "",
		name:    "test",
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	testContent := "test clipboard content"
	err := monitor.SetContent(testContent)
	if err != nil {
		t.Fatalf("SetContent() failed: %v", err)
	}

	// Verify content was set in the mock
	if mockPlatform.setContent != testContent {
		t.Errorf("SetContent() didn't set content correctly in mock: got %q, want %q", mockPlatform.setContent, testContent)
	}
}

func TestMonitor_StartMonitoringWithInterval(t *testing.T) {
	// Use a mock platform for testing that changes content
	mockPlatform := &mockPlatform{
		content:       "initial content",
		name:          "test",
		changeContent: true,
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	// Use a timeout context instead of relying on callback to cancel
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var receivedContent []models.ClipboardContent
	var callCount int

	// Start monitoring with very short interval for testing
	err := monitor.StartMonitoringWithInterval(ctx, 10*time.Millisecond, func(content models.ClipboardContent, m *Monitor) {
		receivedContent = append(receivedContent, content)
		callCount++
	})

	if err != nil && err != context.Canceled {
		t.Fatalf("StartMonitoringWithInterval() failed: %v", err)
	}

	// Wait for the monitoring to complete
	time.Sleep(150 * time.Millisecond)

	if callCount < 1 {
		t.Error("Monitoring callback was not called")
	}

	// Should receive changed content (callback only called when content changes)
	if len(receivedContent) >= 1 {
		if receivedContent[0].Content != "changed content" {
			t.Errorf("First callback content = %q, want %q", receivedContent[0].Content, "changed content")
		}
		if receivedContent[0].Format != "text/plain" {
			t.Errorf("First callback format = %q, want %q", receivedContent[0].Format, "text/plain")
		}
	}
}

func TestMonitor_StartMonitoring(t *testing.T) {
	// Use a mock platform for testing that changes content
	mockPlatform := &mockPlatform{
		content:       "test content",
		name:          "test",
		changeContent: true,
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	// Use a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // Longer timeout for 500ms polling
	defer cancel()

	var receivedContent []models.ClipboardContent
	var callCount int

	// Start monitoring with default interval (which uses clipboardPollInterval = 500ms)
	err := monitor.StartMonitoring(ctx, func(content models.ClipboardContent, m *Monitor) {
		receivedContent = append(receivedContent, content)
		callCount++
	})

	if err != nil && err != context.Canceled {
		t.Fatalf("StartMonitoring() failed: %v", err)
	}

	// Wait for the monitoring to complete
	time.Sleep(2*time.Second + 100*time.Millisecond)

	if callCount < 1 {
		t.Error("Monitoring callback was not called")
	}

	// Should receive changed content (callback only called when content changes)
	if len(receivedContent) >= 1 {
		if receivedContent[0].Content != "changed content" {
			t.Errorf("First callback content = %q, want %q", receivedContent[0].Content, "changed content")
		}
		if receivedContent[0].Format != "text/plain" {
			t.Errorf("First callback format = %q, want %q", receivedContent[0].Format, "text/plain")
		}
	}
}

func TestMonitor_Stop(t *testing.T) {
	// Use a mock platform for testing
	mockPlatform := &mockPlatform{
		content: "stop test content",
		name:    "stop-test",
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	// Stop should not block or panic
	monitor.Stop()

	// Multiple stops should be safe
	monitor.Stop()
	monitor.Stop()
}

func TestMonitor_PlatformInterface(t *testing.T) {
	// Use a mock platform for testing
	mockPlatform := &mockPlatform{
		content: "interface test content",
		name:    "interface-test",
	}

	monitor := &Monitor{
		platform: mockPlatform,
		stopChan: make(chan struct{}),
	}

	// Test GetName
	platformName := monitor.platform.GetName()
	if platformName != "interface-test" {
		t.Errorf("Platform.GetName() = %q, want %q", platformName, "interface-test")
	}

	// Test that we can call platform methods through the monitor
	content, err := monitor.GetCurrentContent()
	if err != nil {
		t.Fatalf("Platform.GetContent() via monitor failed: %v", err)
	}
	if content != "interface test content" {
		t.Errorf("GetCurrentContent() via monitor = %q, want %q", content, "interface test content")
	}
}

// Test platform-specific implementations
func TestWindowsPlatform(t *testing.T) {
	// This test would normally be skipped on non-Windows systems
	// but we can at least verify the constants and structure
	if cfUnicodeText != 13 {
		t.Errorf("cfUnicodeText constant = %d, want 13", cfUnicodeText)
	}
	if gmemMoveable != 0x0002 {
		t.Errorf("gmemMoveable constant = %d, want 0x0002", gmemMoveable)
	}
}

// Note: Platform-specific types (linuxPlatform, darwinPlatform) are only available
// when building on their respective platforms due to build tags.
// These tests would be run on the appropriate platform during CI.

// Integration test with actual platform (be careful with this)
func TestMonitor_Integration(t *testing.T) {
	t.Skip("Skipping integration test as it requires actual clipboard access")
}
