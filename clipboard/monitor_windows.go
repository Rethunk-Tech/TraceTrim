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
	defer closeClipboard.Call()

	// Get data in Unicode text format
	hMem, _, _ := getClipboardData.Call(CF_UNICODETEXT)
	if hMem == 0 {
		return "", fmt.Errorf("no text data in clipboard")
	}

	// Lock the memory and get pointer
	lockRet, _, _ := globalLock.Call(hMem)
	if lockRet == 0 {
		return "", fmt.Errorf("failed to lock global memory")
	}
	defer globalUnlock.Call(hMem)

	// Convert UTF-16 bytes to Go string
	utf16Ptr := (*uint16)(unsafe.Pointer(lockRet))
	length := 0
	for *utf16Ptr != 0 {
		utf16Ptr = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(utf16Ptr)) + 2))
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
	size := len(utf16) * 2

	// Open clipboard
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to open clipboard")
	}
	defer closeClipboard.Call()

	// Empty clipboard
	emptyClipboard.Call()

	// Allocate global memory
	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(size+2))
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
	globalUnlock.Call(hMem)

	// Set clipboard data
	setClipboardDataRet, _, _ := setClipboardData.Call(CF_UNICODETEXT, hMem)
	if setClipboardDataRet == 0 {
		globalFree.Call(hMem)
		return fmt.Errorf("failed to set clipboard data")
	}

	return nil
}

// Windows API constants and function declarations
var (
	user32           = syscall.NewLazyDLL("user32.dll")
	kernel32         = syscall.NewLazyDLL("kernel32.dll")

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
	CF_UNICODETEXT = 13
	GMEM_MOVEABLE  = 0x0002
)
