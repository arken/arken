package config

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// DiskInfo is a struct that provides information about the mounted directory size.
type DiskInfo struct {
	AvailableBytes uint64
	PoolSizeBytes  uint64
}

// DiskInfoProvider is an interface is what makes this part of the program cross
// platform. The Unix and windows version of DiskInfo implement this method interface.
type DiskInfoProvider interface {
	Init()
	Refresh()
	GetAvailableBytes() uint64
	SetAvailableBytes(uint64)
	GetPoolSizeBytes() uint64
	SetPoolSizeBytes(uint64)
	GetPrettyPoolSize() string
}

// GetDiskInfo returns the disk info struct.
func (di *DiskInfo) GetDiskInfo() *DiskInfo {
	return di
}

// GetPoolSizeBytes returns the size of the storage pool.
func (di *DiskInfo) GetPoolSizeBytes() uint64 {
	return di.PoolSizeBytes
}

// SetPoolSizeBytes sets the size of the storage pool.
func (di *DiskInfo) SetPoolSizeBytes(new uint64) {
	di.PoolSizeBytes = new
}

// GetAvailableBytes gets the number of bytes still free on the drive.
func (di *DiskInfo) GetAvailableBytes() uint64 {
	return di.AvailableBytes
}

// SetAvailableBytes sets the number of free bytes on the drive.
func (di *DiskInfo) SetAvailableBytes(new uint64) {
	di.AvailableBytes = new
}

// GetPrettyPoolSize outputs a pretty string of the pool size.
func (di *DiskInfo) GetPrettyPoolSize() string {
	if di.PoolSizeBytes < 1000000000000 { //less than a TB
		return fmt.Sprintf("%vGB", toUnit(di.PoolSizeBytes, 9))
	}
	return fmt.Sprintf("%vTB", toUnit(di.PoolSizeBytes, 12))
}

// GlobalDiskInfo is the Global Configuration struct for Arken disk stats
var GlobalDiskInfo DiskInfo

//ParsePoolSize parses the string contained in the global config as the max pool size,
//and stores the value it comes up with in bytes in the struct GlobalDiskInfo. If
//for some reason the input cannot be understood or the user attempts to allocate
//more storage than they have available, a default value of the user's capacity
//minus 10 GB will be put in place. This means that the minimum amount of space
//a user must have available is 10 GB.
//The string must be in the following format and order:
//  1. A base-10 number, can be floating point
//  2. One of the following: B, KB, MB, GB, or TB. It's case sensitive to avoid
//     bit/byte confusion.
//There can be any amount of whitespace before and after either of the elements.
//  "3000MB", "  10GB   ", "10 GB", "1.75TB", ".5 TB" will all work.
//  "1.TB", "10gb", "0xfa5MB" will not work
func ParsePoolSize(dip DiskInfoProvider) {
	dip.Refresh()
	max := Global.General.PoolSize
	defaultSizeB := int64(dip.GetAvailableBytes()) - 10000000000 //available - 10GB
	defaultSizeGB := toUnit(uint64(defaultSizeB), 9)             //for use in strings
	if defaultSizeB < 0 {                                        //user has less than 10GB available
		log.Fatal("Not enough free storage on this device, 10GB or more is required")
	}
	parentRegex := regexp.MustCompile(
		"\\A\\s*([0-9]*\\.)?([0-9]\\d*)\\s*[KMGT]?B\\s*$",
	)
	if strings.EqualFold(max, "max") { //case insensitive comparison
		dip.SetPoolSizeBytes(dip.GetAvailableBytes())
		Global.General.PoolSize = dip.GetPrettyPoolSize()
	} else if parentRegex.MatchString(max) {
		dip.SetPoolSizeBytes(ParseWellFormedPoolSize(max))
	} else { //did not match parent regex
		log.Printf("Unable to understand \"%v\" as max pool size,"+
			" using %v GB instead\n", max, defaultSizeGB)
		log.Println(`
The max pool size string must be in the following format and order:
    1. A base-10 number, can be floating point
    2. One of the following: B, KB, MB, GB, or TB. It's case sensitive to avoid
     bit/byte confusion.
There can be any amount of whitespace before and after either of the elements.
"3000MB", "  10GB   ", "10 GB", "1.75TB", ".5 TB" will all work.
"1.TB", "10gb", "0xfa5MB" will not work`)
		dip.SetPoolSizeBytes(uint64(defaultSizeB))
		Global.General.PoolSize = dip.GetPrettyPoolSize()
	}
	if dip.GetPoolSizeBytes() > dip.GetAvailableBytes() {
		log.Printf("Less than the requested %vGB are free on this computer, "+
			"using %vGB instead\n", toUnit(dip.GetPoolSizeBytes(), 9), defaultSizeGB)
		dip.SetPoolSizeBytes(uint64(defaultSizeB))
		Global.General.PoolSize = dip.GetPrettyPoolSize()
	}
	printResults(dip)
}

// ParseWellFormedPoolSize parses a string that passed the regex test in ParsePoolSize(). It extracts the
// number and unit, returning the number of bytes indicated by the string.
// IE: parseWellFormedPoolSize("10GB") = 10,000,000,000
func ParseWellFormedPoolSize(str string) uint64 {
	//extract the number
	bytesStr := regexp.MustCompile("([0-9]*\\.)?([0-9]\\d*)").FindString(str)
	//extract the unit of storage
	unitsStr := regexp.MustCompile("[KMGT]?B").FindString(str)
	bytesFloat, _ := strconv.ParseFloat(bytesStr, 64)
	switch unitsStr {
	case "TB":
		bytesFloat *= math.Pow10(12)
	case "GB":
		bytesFloat *= math.Pow10(9)
	case "MB":
		bytesFloat *= math.Pow10(6)
	case "KB":
		bytesFloat *= math.Pow10(3)
	}
	bytes := uint64(bytesFloat)
	if bytes < 1000000000 {
		log.Fatal("Arken requires an allocation of at least 1GB")
	}
	return bytes
}

//Takes in a uint64 of bytes, return a float64 representing the amount of bytes
//in gigabytes.
func toUnit(bytes uint64, pow int) float64 {
	return float64(bytes) / math.Pow10(pow)
}

//This function prints the final results of the parsing of the pool size. It
//attempts to print the detected storage and storage allocated to Arken in a
//readable unit.
func printResults(dip DiskInfoProvider) {
	poolStr := BytesToUnitString(dip.GetPoolSizeBytes())
	availStr := BytesToUnitString(dip.GetAvailableBytes())
	log.Printf("Detected %v of storage available on this "+
		"device, using %v (0x%x bytes)\n", availStr, poolStr, dip.GetPoolSizeBytes())
}

// BytesToUnitString Given a number of bytes, Returns a string that represents the number of bytes in
// a sensible unit.
func BytesToUnitString(bytes uint64) string {
	var pow int
	var unit string
	if bytes >= 1000000000000 {
		pow = 12
		unit = "TB"
	} else if bytes >= 1000000000 {
		pow = 9
		unit = "GB"
	} else if bytes >= 1000000 {
		pow = 6
		unit = "MB"
	} else if bytes >= 1000 {
		pow = 3
		unit = "KB"
	} else {
		pow = 0
		unit = "B"
	}
	return fmt.Sprintf("%v%v", toUnit(bytes, pow), unit)
}

//TODO write a generic UnitStringToBytes() that does the inverse of ^
