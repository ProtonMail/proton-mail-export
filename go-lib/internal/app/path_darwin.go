//go:build darwin

package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>
#import <stdio.h>
char const* getMacOSDownloadsDir() {
    @autoreleasepool{
        NSArray *paths = NSSearchPathForDirectoriesInDomains(NSDownloadsDirectory, NSUserDomainMask, YES);
		// The buffer returned by UTF8String's lifetime does not exceed the lifetime of paths, so we must copy it.
		return strdup([[paths objectAtIndex:0] UTF8String]);// memory allocated by strudp is to be related be freed on the Go side.
	}
}
*/
import "C"
import "unsafe"

func getDownloadDir() string {
	var cStr *C.char = C.getMacOSDownloadsDir()
	var result string = C.GoString(cStr)
	C.free(unsafe.Pointer(cStr))
	return result
}
