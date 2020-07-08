// +build dragonfly freebsd hurd illumos linux netbsd openbsd solaris

package config

import (
	"golang.org/x/sys/unix"
	"log"
)

var GlobalDiskInfo DiskInfo

type DiskInfo struct {
	AvailableBytes uint64
	PoolSizeBytes  uint64
	isUnix         bool
}

//Initiates a DiskInfo struct, checking what os the program is running on and
//then calling the appropriate syscall to get the requisite information about
//the disk.
func (di* DiskInfo) Init() {
	fs := unix.Statfs_t{}
	err := unix.Statfs(".", &fs)
	if err != nil {
		log.Fatal(err)
	}
	di.AvailableBytes = fs.Bavail * uint64(fs.Bsize)
}

//Refreshes the info with a new syscall. This is not called in GetAvailableBytes()
//because syscalls are expensive.
func (di* DiskInfo) Refresh() {
	di.Init()
}
