package config

import (
    "golang.org/x/sys/unix"
    "log"
    "math"
    "regexp"
    "runtime"
    "strconv"
    "strings"
)

var diskInfo DiskInfo

type DiskInfo struct {
    AvailableBytes uint64
    ArkenAllocatedBytes uint64
    isUnix bool
}

//Initiates a DiskInfo struct, checking what os the program is running on and
//then calling the appropriate syscall to get the requisite information about
//the disk.
func (di* DiskInfo) Init() {
    var unixes = []string{ "dragonfly", "freebsd", "hurd", "illumos",
        "linux", "netbsd", "openbsd", "solaris" }
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
            di.AvailableBytes = fs.Bavail * uint64(fs.Bsize)
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
        di.AvailableBytes = fs.Bavail * uint64(fs.Bsize)
    } else {
        //TODO: figure out what syscall to use on windows.
    }
}

func ParseUserDiskInput(di* DiskInfo) {
    di.Refresh()
    max := Global.General.PoolSize
    defaultSizeB := int64(di.AvailableBytes) - 10000000000 //available - 10GB
    defaultSizeGB := float64(defaultSizeB) / math.Pow10(9) //for use in strings
    //						v 1 GB v
    if di.AvailableBytes < 1000000000 || defaultSizeB < 0 {
        log.Fatal("Not enough free storage on this device")
    }
    var poolSizeB uint64
    var err error
    parentRegex := regexp.MustCompile("\\A\\s*[1-9]\\d*\\s*[MGT]B\\s*$")
    if strings.EqualFold(max, "max") { //case insensitive comparison
        poolSizeB = di.AvailableBytes
    } else if parentRegex.MatchString(max) {
        bytesStr := regexp.MustCompile("[1-9]\\d*").FindString(max) //extract the number
        unitStr := regexp.MustCompile("[MGT]B").FindString(max) //extract the unit (GB/MB)
        poolSizeB, err = strconv.ParseUint(bytesStr, 10, 64) //convert number to uint64
        if err != nil { //theoretically, the regex should avoid parsing errors.
            log.Printf("Unable to understand \"%v\" as max pool size," +
                " using %v GB instead\n", max, defaultSizeGB)
            poolSizeB = uint64(defaultSizeB)
        } else {
            poolSizeB = parseStorageUnit(&unitStr, poolSizeB)
            if poolSizeB < 1000000000 {
                log.Fatal("Arken requires an allocation of at least 1 GB")
            }
        }
    } else { //did not match parent regex
        log.Printf("Unable to understand \"%v\" as max pool size," +
            " attempting to use %v GB instead\n", max, defaultSizeGB)
        poolSizeB = uint64(defaultSizeB)
    }
    if poolSizeB > di.AvailableBytes {
        log.Printf("Less than the requested %v GB is/are free on this computer, " +
            "using %v GB instead\n", float64(poolSizeB) / math.Pow10(9), defaultSizeGB)
        poolSizeB = uint64(defaultSizeB)
    }
    diskInfo.ArkenAllocatedBytes = poolSizeB
    log.Printf("Using %v GB of storage", float64(poolSizeB) / math.Pow10(9))
}

func parseStorageUnit(unitStr *string, num uint64) uint64 {
    result := num
    if *unitStr == "TB" {
        result *= uint64(math.Pow10(12))
    } else if *unitStr == "GB" {
        result *= uint64(math.Pow10(9))
    } else {
        result *= uint64(math.Pow10(6))
    }
    return result
}
