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
	GetDiskInfo() *DiskInfo
}

// GetDiskInfo returns the disk info struct.
func (di *DiskInfo) GetDiskInfo() *DiskInfo {
	return di
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
	di := dip.GetDiskInfo()
	di.Refresh()
	max := Global.General.PoolSize
	defaultSizeB := int64(di.AvailableBytes) - 10000000000 //available - 10GB
	defaultSizeGB := toGB(uint64(defaultSizeB))            //for use in strings
	if defaultSizeB < 0 {                                  //user has less than 10GB available
		log.Fatal("Not enough free storage on this device, 10 GB or more is required")
	}
	var poolSizeB uint64
	parentRegex := regexp.MustCompile(
		"\\A\\s*([0-9]*\\.)?([0-9]\\d*)\\s*[KMGT]?B\\s*$",
	)
	if strings.EqualFold(max, "max") { //case insensitive comparison
		poolSizeB = di.AvailableBytes
		Global.General.PoolSize = fmt.Sprintf("%vB", di.AvailableBytes)
	} else if parentRegex.MatchString(max) {
		poolSizeB = parseWellFormedPoolSize(max)
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
		poolSizeB = uint64(defaultSizeB)
		Global.General.PoolSize = fmt.Sprintf("%vB", defaultSizeB)
	}
	if poolSizeB > di.AvailableBytes {
		log.Printf("Less than the requested %v GB are free on this computer, "+
			"using %v GB instead\n", toGB(poolSizeB), defaultSizeGB)
		poolSizeB = uint64(defaultSizeB)
		Global.General.PoolSize = fmt.Sprintf("%vB", defaultSizeB)
	}
	di.PoolSizeBytes = poolSizeB
	printResults(dip)
}

//parses a string that passed the regex test in ParsePoolSize(). It extracts the
//number and unit, returning the number of bytes indicated by the string.
//IE: parseWellFormedPoolSize("10GB") = 10,000,000,000
func parseWellFormedPoolSize(str string) uint64 {
	//extract the number
	bytesStr := regexp.MustCompile("([0-9]*\\.)?([0-9]\\d*)").FindString(str)
	//extract the unit of storage
	unitStr := regexp.MustCompile("[KMGT]?B").FindString(str)
	bytesFloat, _ := strconv.ParseFloat(bytesStr, 64)
	switch unitStr {
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
		log.Fatal("Arken requires an allocation of at least 1 GB")
	}
	return bytes
}

//Takes in a uint64 of bytes, return a float64 representing the amount of bytes
//in gigabytes.
func toGB(bytes uint64) float64 {
	return float64(bytes) / math.Pow10(9)
}

//This function prints the final results of the parsing of the pool size. It
//attempts to print the detected storage and storage allocated to Arken in a
//readable unit.
func printResults(dip DiskInfoProvider) {
	di := dip.GetDiskInfo()
	availPow10 := 0
	poolPow10 := 0
	for i := 0; i <= 12; i += 3 {
		availLen := numLen(di.AvailableBytes / uint64(math.Pow10(i)))
		if availLen <= 3 && availLen > 0 {
			availPow10 = i
		}
		poolLen := numLen(di.PoolSizeBytes / uint64(math.Pow10(i)))
		if poolLen <= 3 && poolLen > 0 {
			poolPow10 = i
		}
	}
	available := float64(di.AvailableBytes) / math.Pow10(availPow10)
	pool := float64(di.PoolSizeBytes) / math.Pow10(poolPow10)
	log.Printf("Detected about %.2f %v of storage available on this " +
		"device, using %v %v (0x%x bytes)\n",
		available, getUnit(availPow10), pool, getUnit(poolPow10), di.PoolSizeBytes)
}

//Returns the appropriate unit of storage for the given power of 10. B up to TB.
func getUnit(pow10 int) string {
	switch pow10 {
	case 12:
		return "TB"
	case 9:
		return "GB"
	case 6:
		return "MB"
	case 3:
		return "KB"
	case 0:
		return "B"
	default:
		return ""
	}
}

//simply returns the number of digits in a given number.
func numLen(num uint64) int {
	if num == 0 {
		return 0
	}
	return 1 + numLen(num / 10)
}
