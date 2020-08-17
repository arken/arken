// +build windows

package config

import (
	"log"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Init cretes a DiskInfo with methods that make unix system calls.
func (di *DiskInfo) Init() {
	os.MkdirAll(filepath.Dir(Global.Sources.Storage), os.ModePerm)

	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	var freeBytesAvailableToCaller uint64
	_, _, err := c.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(Global.Sources.Storage))),
		uintptr(unsafe.Pointer(&freeBytesAvailableToCaller)),
		uintptr(0), //don't care about this value
		uintptr(0), //don't care about this value
	)
	//the syscall returns an error that is always non-nil, and if it worked,
	//this is the error it will return, so this is like != nil in this
	//particular case. I check for nil just in case.
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(err)
	}
	di.AvailableBytes = freeBytesAvailableToCaller
}

//Refreshes the info with a new syscall.
func (di *DiskInfo) Refresh() {
	di.Init()
}
