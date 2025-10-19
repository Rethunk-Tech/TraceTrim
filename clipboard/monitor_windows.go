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
	return &standardPlatform{}, nil
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
		return "", fmt.Errorf("failed to open Windows clipboard: access denied or clipboard in use")
	}
	defer func() {
		closeClipboard.Call() //nolint:errcheck // Ignore errors in defer
	}()

	// Get data in Unicode text format
	hMem, _, _ := getClipboardData.Call(cfUnicodeText)
	if hMem == 0 {
		return "", fmt.Errorf("no Unicode text data available in Windows clipboard (format: %d)", cfUnicodeText)
	}

	// Lock the memory and get pointer
	lockRet, _, _ := globalLock.Call(hMem)
	if lockRet == 0 {
		return "", fmt.Errorf("failed to lock Windows clipboard memory object")
	}

	// Ensure memory is always unlocked when we're done
	defer func() {
		if lockRet != 0 {
			globalUnlock.Call(hMem) //nolint:errcheck // Ignore errors in defer
		}
	}()

	// Get the size of the clipboard data
	size, _, _ := globalSize.Call(hMem)
	if size == 0 {
		return "", fmt.Errorf("failed to get clipboard data size")
	}

	// Calculate the length in uint16 units (size includes null terminator)
	length := int(size) / int(utf16CharSize)

	// Copy data to a buffer using Windows API pattern
	utf16Slice := make([]uint16, length)
	srcPtr := uintptr(lockRet)
	dstPtr := uintptr(unsafe.Pointer(&utf16Slice[0]))

	// Use Windows memory copy function for safety
	kernel32.NewProc("RtlMoveMemory").Call(dstPtr, srcPtr, size)

	// Convert to Go string
	goStr := syscall.UTF16ToString(utf16Slice)

	return goStr, nil
}

// SetContent sets text content to Windows clipboard
func (w *windowsPlatform) SetContent(content string) error {
	// Convert string to UTF-16 bytes (includes null terminator)
	utf16, err := syscall.UTF16FromString(content)
	if err != nil {
		return fmt.Errorf("failed to convert string to UTF-16: %w", err)
	}

	// Calculate size needed for the memory block (includes null terminator)
	size := len(utf16) * utf16CharSize

	// Open clipboard
	ret, _, _ := openClipboard.Call(0)
	if ret == 0 {
		return fmt.Errorf("failed to open Windows clipboard: access denied or clipboard in use")
	}
	defer func() {
		closeClipboard.Call() //nolint:errcheck // Ignore errors in defer
	}()

	// Empty clipboard
	ret, _, _ = emptyClipboard.Call()
	if ret == 0 {
		return fmt.Errorf("failed to empty Windows clipboard")
	}

	// Allocate global memory (moveable and zero-initialized)
	hMem, _, _ := globalAlloc.Call(gmemMoveable|gmemZeroInit, uintptr(size))
	if hMem == 0 {
		return fmt.Errorf("failed to allocate Windows global memory: insufficient memory")
	}

	// Lock memory and copy data
	lockRet, _, _ := globalLock.Call(hMem)
	if lockRet == 0 {
		globalFree.Call(hMem) //nolint:errcheck // Ignore errors in cleanup
		return fmt.Errorf("failed to lock Windows clipboard memory object")
	}

	// Copy UTF-16 data including null terminator using Windows API pattern
	srcPtr := uintptr(unsafe.Pointer(&utf16[0]))
	dstPtr := uintptr(lockRet)

	// Use Windows memory copy function for safety - copy all data including null terminator
	kernel32.NewProc("RtlMoveMemory").Call(dstPtr, srcPtr, uintptr(size))

	// Unlock memory - GlobalUnlock returns 0 on failure, non-zero on success
	unlockRet, _, _ := globalUnlock.Call(hMem)
	if unlockRet == 0 {
		globalFree.Call(hMem) //nolint:errcheck // Ignore errors in cleanup
		return fmt.Errorf("failed to unlock Windows clipboard memory object")
	}

	// Set clipboard data (Windows will own the memory after this call)
	setClipboardDataRet, _, _ := setClipboardData.Call(cfUnicodeText, hMem)
	if setClipboardDataRet == 0 {
		// If SetClipboardData fails, we need to free the memory since Windows doesn't own it
		globalFree.Call(hMem) //nolint:errcheck // Ignore errors in cleanup
		return fmt.Errorf("failed to set Unicode text data in Windows clipboard")
	}

	// Memory is now owned by Windows clipboard - don't free it

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
	globalSize       = kernel32.NewProc("GlobalSize")
)

const (
	cfUnicodeText = 13 // CF_UNICODETEXT
	gmemMoveable  = 0x0002
	gmemZeroInit  = 0x0040
)
