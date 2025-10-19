//go:build linux
// +build linux

package clipboard

import (
	"fmt"
	"os/exec"
	"strings"
)

// linuxPlatform implements Platform interface for Linux
type linuxPlatform struct{}

// getPlatform returns the appropriate platform implementation for Linux
func getPlatform() (Platform, error) {
	return &linuxPlatform{}, nil
}

// GetName returns the platform name
func (l *linuxPlatform) GetName() string {
	return "Linux"
}

// GetContent retrieves text content from Linux clipboard
func (l *linuxPlatform) GetContent() (string, error) {
	// Try xclip first (supports both X11 and Wayland via XWayland)
	content, err := l.getContentWithXclip()
	if err == nil {
		return content, nil
	}

	// Fall back to xsel
	content, err = l.getContentWithXsel()
	if err != nil {
		return "", fmt.Errorf("failed to get clipboard content (tried xclip and xsel): %w", err)
	}

	return content, nil
}

// getContentWithXclip retrieves clipboard content using xclip
func (l *linuxPlatform) getContentWithXclip() (string, error) {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("xclip failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// getContentWithXsel retrieves clipboard content using xsel
func (l *linuxPlatform) getContentWithXsel() (string, error) {
	cmd := exec.Command("xsel", "-ob")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("xsel failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// SetContent sets text content to Linux clipboard
func (l *linuxPlatform) SetContent(content string) error {
	// Try xclip first
	err := l.setContentWithXclip(content)
	if err == nil {
		return nil
	}

	// Fall back to xsel
	err = l.setContentWithXsel(content)
	if err != nil {
		return fmt.Errorf("failed to set clipboard content (tried xclip and xsel): %w", err)
	}

	return nil
}

// setContentWithXclip sets clipboard content using xclip
func (l *linuxPlatform) setContentWithXclip(content string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard", "-i")
	cmd.Stdin = strings.NewReader(content)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("xclip set failed: %w", err)
	}

	return nil
}

// setContentWithXsel sets clipboard content using xsel
func (l *linuxPlatform) setContentWithXsel(content string) error {
	cmd := exec.Command("xsel", "-ib")
	cmd.Stdin = strings.NewReader(content)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("xsel set failed: %w", err)
	}

	return nil
}
