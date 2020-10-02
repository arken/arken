package ipfs

import (
	"log"
	"time"
)

var (
	output chan uint64
	signal chan bool
	delta  chan int64
)

func init() {
	signal = make(chan bool)
	output = make(chan uint64)
	delta = make(chan int64)

	go getRepoSizeDaemon()
}

func getRepoSizeDaemon() {
	var (
		err       error
		usedSpace uint64
	)
	checkTime := time.Now().AddDate(0, -1, 0)

	for {
		select {
		case change := <-delta:
			usedSpace = uint64(int64(usedSpace) + change)
			continue

		case <-signal:
			if time.Since(checkTime) >= time.Hour {
				usedSpace, err = calculateRepoSize()
				if err != nil {
					log.Fatal(err)
				}
				checkTime = time.Now()
			}
			output <- usedSpace
			continue
		}
	}
}

// AdjustRepoSize adds or removes the given amount from the calculated
// IPFS store size.
func AdjustRepoSize(change int64) {
	delta <- change
}

// GetRepoSize returns the size of the repo in bytes.
func GetRepoSize() (result uint64, err error) {
	signal <- true
	return <-output, nil
}

// calculateRepoSize returns the size of the repo in bytes.
func calculateRepoSize() (result uint64, err error) {
	out, err := node.Repo.GetStorageUsage()
	if err != nil {
		return result, err
	}
	return out, nil
}
