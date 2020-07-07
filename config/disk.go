package config

import (
    "golang.org/x/sys/unix"
    "log"
    "runtime"
)

type DiskInfo struct {
    availableBytes uint64
    isUnix bool
}

//Initiates a DiskInfo struct, checking what os the program is running on and
//then calling the appropriate syscall to get the requisite information about
//the disk.
func (di* DiskInfo) Init() {
    var unixes = []string{ "dragonfly", "freebsd", "hurd", "illumos",
        "linux", "netbsd", "openbsd", "solaris"}
    if runtime.GOOS == "windows" {
        di.isUnix = false
        //TODO: figure out what syscall to use on windows.
    } else  {
        for _, os := range unixes {
            if runtime.GOOS == os {
                di.isUnix = true
                break
            }
        }
        if di.isUnix {
            fs := unix.Statfs_t{}
            err := unix.Statfs(".", &fs)
            if err != nil {
                log.Fatal(err)
            }
            di.availableBytes = fs.Bavail * uint64(fs.Bsize)
        } else {
            panic("Unrecognized operating system \"" + runtime.GOOS + "\".")
        }
    }
}

//Refresh the info with a new syscall. This is not called in GetAvailableBytes()
//because syscalls are expensive.
func (di* DiskInfo) Refresh() {
    if di.isUnix {
        fs := unix.Statfs_t{}
        err := unix.Statfs(".", &fs)
        if err != nil {
            log.Fatal(err)
        }
        di.availableBytes = fs.Bavail * uint64(fs.Bsize)
    } else {
        //TODO: figure out what syscall to use on windows.
    }
}

//Returns the amount of available space on disk in bytes. This is not necessarily
//an up-to-date value, and calling this function will not result in a syscall. This
//will return the value established either on Init() or Refresh()
func (di* DiskInfo) GetAvailableBytes() uint64 {
    return di.availableBytes
}

func ParseUserDiskInput(capacity string) {
    //Here I will eventually check how much space the user wants to allow to be
    //allocated.
}

