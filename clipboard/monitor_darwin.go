//go:build darwin
// +build darwin

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

char* getClipboardText() {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSString *bestType = [pasteboard availableTypeFromArray:[NSArray arrayWithObject:NSPasteboardTypeString]];
        if (!bestType) {
            return NULL;
        }

        NSString *string = [pasteboard stringForType:NSPasteboardTypeString];
        if (!string) {
            return NULL;
        }

        const char *utf8String = [string UTF8String];
        if (!utf8String) {
            return NULL;
        }

        // Allocate memory for the C string
        char *result = strdup(utf8String);
        return result;
    }
}

bool setClipboardText(const char *text) {
    @autoreleasepool {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];

        NSString *nsString = [NSString stringWithUTF8String:text];
        if (!nsString) {
            return false;
        }

        BOOL success = [pasteboard writeObjects:[NSArray arrayWithObject:nsString]];
        return success ? true : false;
    }
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// getPlatform returns the appropriate platform implementation for macOS
func getPlatform() (Platform, error) {
	return &darwinPlatform{}, nil
}

// darwinPlatform implements Platform interface for macOS
type darwinPlatform struct{}

// GetName returns the platform name
func (d *darwinPlatform) GetName() string {
	return "macOS"
}

// GetContent retrieves text content from macOS clipboard
func (d *darwinPlatform) GetContent() (string, error) {
	cStr := C.getClipboardText()
	if cStr == nil {
		return "", fmt.Errorf("no text data in clipboard or failed to read")
	}
	defer C.free(unsafe.Pointer(cStr))

	return C.GoString(cStr), nil
}

// SetContent sets text content to macOS clipboard
func (d *darwinPlatform) SetContent(content string) error {
	cStr := C.CString(content)
	defer C.free(unsafe.Pointer(cStr))

	success := C.setClipboardText(cStr)
	if !success {
		return fmt.Errorf("failed to set clipboard content")
	}

	return nil
}
