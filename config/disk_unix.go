// +build dragonfly freebsd hurd illumos linux netbsd openbsd solaris

package config

import (
	"golang.org/x/sys/unix"
	"log"
)

//Initiates a DiskInfo with methods that make unix system calls.
func (di* DiskInfo) Init() {
	fs := unix.Statfs_t{}
	err := unix.Statfs(".", &fs)
	if err != nil {
		log.Fatal(err)
	}
	di.AvailableBytes = fs.Bavail * uint64(fs.Bsize)
}

//Refreshes the info with a new syscall.
func (di* DiskInfo) Refresh() {
	di.Init()
}
