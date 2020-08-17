// +build dragonfly freebsd hurd illumos linux netbsd openbsd solaris darwin

package config

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Init initializes a DiskInfo with methods that make unix system calls.
func (di *DiskInfo) Init() {
	os.MkdirAll(filepath.Dir(Global.Sources.Storage), os.ModePerm)

	fs := unix.Statfs_t{}
	err := unix.Statfs(Global.Sources.Storage, &fs)
	if err != nil {
		log.Fatal(err)
	}
	di.AvailableBytes = fs.Bavail * uint64(fs.Bsize)
}

// Refresh updates the info with a new syscall.
func (di *DiskInfo) Refresh() {
	di.Init()
}
