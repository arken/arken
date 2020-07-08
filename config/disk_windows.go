// +build windows

package config

import (
	"golang.org/x/sys/windows"
	"log"
	"os"
	"unsafe"
)

//Initiates a DiskInfo struct, checking what os the program is running on and
//then calling the appropriate syscall to get the requisite information about
//the disk.
func (di* DiskInfo) Init() {
	wd, _ := os.Getwd() //working directory
	h := windows.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")
	var freeBytesAvailableToCaller uint64
	_, _, err := c.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(wd))),
		uintptr(unsafe.Pointer(&freeBytesAvailableToCaller)),
		uintptr(0), //don't care about these values
		uintptr(0), //don't care about these values
	)
	//the syscall returns an error that is always non-nil, and if it worked,
	//this is the error it will return, so this is like != nil in this
	//particular case. I check for nil just in case
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(err)
	}
	di.AvailableBytes = freeBytesAvailableToCaller
}

//Refreshes the info with a new syscall.
func (di* DiskInfo) Refresh() {
	di.Init()
}