//go:build windows
// +build windows

package clipboard

import (
	"fmt"
	"syscall"
	"unsafe"
)

// windowsPlatform implements Platform interface for Windows
type windowsPlatform struct{}

const (
	// UTF-16 character size in bytes
	utf16CharSize = 2
)

// getPlatform returns the appropriate platform implementation
func getPlatform() (Platform, error) {
	return &windowsPlatform{}, nil
}

// GetName returns the platform name
func (w *windowsPlatform) GetName() string {
	return "Windows"
}

// GetContent retrieves text content from Windows clipboard
func (w *windowsPlatform) GetContent() (string, error) {
	// Open clipboard
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return "", fmt.Errorf("failed to open clipboard")
	}
	defer func() {
		if _, _, err := closeClipboard.Call(); err != nil {
			// Log error but don't fail the operation
		}
	}()

	// Get data in Unicode text format
	hMem, _, _ := getClipboardData.Call(cfUnicodeText)
	if hMem == 0 {
		return "", fmt.Errorf("no text data in clipboard")
	}

	// Lock the memory and get pointer
	lockRet, _, _ := globalLock.Call(hMem)
	if lockRet == 0 {
		return "", fmt.Errorf("failed to lock global memory")
	}
	defer func() {
		if _, _, err := globalUnlock.Call(hMem); err != nil {
			// Log error but don't fail the operation
		}
	}()

	// Convert UTF-16 bytes to Go string
	utf16Ptr := (*uint16)(unsafe.Pointer(lockRet))
	length := 0
	for *utf16Ptr != 0 {
		utf16Ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(utf16Ptr)) + utf16CharSize))
		length++
	}

	utf16Ptr = (*uint16)(unsafe.Pointer(lockRet))
	goStr := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(utf16Ptr))[:length:length])

	return goStr, nil
}

// SetContent sets text content to Windows clipboard
func (w *windowsPlatform) SetContent(content string) error {
	// Convert string to UTF-16 bytes
	utf16, err := syscall.UTF16FromString(content)
	if err != nil {
		return fmt.Errorf("failed to convert string to UTF-16: %w", err)
	}

	// Calculate size needed for the memory block
	size := len(utf16) * utf16CharSize

	// Open clipboard
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to open clipboard")
	}
	defer func() {
		if _, _, err := closeClipboard.Call(); err != nil {
			// Log error but don't fail the operation
		}
	}()

	// Empty clipboard
	if _, _, err := emptyClipboard.Call(); err != nil {
		return fmt.Errorf("failed to empty clipboard")
	}

	// Allocate global memory
	hMem, _, _ := globalAlloc.Call(gmemMoveable, uintptr(size+utf16CharSize))
	if hMem == 0 {
		return fmt.Errorf("failed to allocate global memory")
	}

	// Lock memory and copy data
	lockRet, _, _ := globalLock.Call(hMem)
	if lockRet == 0 {
		globalFree.Call(hMem)
		return fmt.Errorf("failed to lock global memory")
	}

	// Copy UTF-16 data
	dest := (*[1 << 20]uint16)(unsafe.Pointer(lockRet))
	copy((*[1 << 20]uint16)(dest)[:len(utf16):len(utf16)], utf16)

	// Unlock memory
	if _, _, err := globalUnlock.Call(hMem); err != nil {
		globalFree.Call(hMem)
		return fmt.Errorf("failed to unlock global memory")
	}

	// Set clipboard data
	setClipboardDataRet, _, _ := setClipboardData.Call(cfUnicodeText, hMem)
	if setClipboardDataRet == 0 {
		globalFree.Call(hMem)
		return fmt.Errorf("failed to set clipboard data")
	}

	return nil
}

// Windows API constants and function declarations
var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	openClipboard    = user32.NewProc("OpenClipboard")
	closeClipboard   = user32.NewProc("CloseClipboard")
	getClipboardData = user32.NewProc("GetClipboardData")
	setClipboardData = user32.NewProc("SetClipboardData")
	emptyClipboard   = user32.NewProc("EmptyClipboard")
	globalAlloc      = kernel32.NewProc("GlobalAlloc")
	globalFree       = kernel32.NewProc("GlobalFree")
	globalLock       = kernel32.NewProc("GlobalLock")
	globalUnlock     = kernel32.NewProc("GlobalUnlock")
)

const (
	cfUnicodeText = 13
	gmemMoveable  = 0x0002
)
