package engine

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

// NetworkLimit is true if the node has hit it's download limit for the month.
var NetworkLimit bool

// Run manages balancing new and at risk files
// between nodes.
func Run(new, remotes, output chan database.FileKey) (err error) {
	keysets := make(map[string]int)
	input := make(chan database.FileKey, 10)

	for set := range config.Keysets {
		name := strings.Split(filepath.Base(config.Keysets[set].URL), ".")[0]
		keysets[name] = config.Keysets[set].Replications

		// Pin Lighthouse File to determine the size of the active cluster.
		err = ipfs.Pin(config.Keysets[set].LightHouseFileID)
		if err != nil {
			return err
		}
	}

	// Determine the possible number of threads for the system's CPU
	workers := genNumWorkers()
	// Generate Worker Threads
	for i := 0; i < workers; i++ {
		go runWorker(keysets, input, output, i)
	}

	for {
		fmt.Printf("New: %d, Remotes: %d, Output: %d\n", len(new), len(remotes), len(output))
		if NetworkLimit {
			select {
			case entry := <-new:
				output <- entry
				continue
			case entry := <-remotes:
				output <- entry
				continue
			default:
				time.Sleep(15 * time.Second)
			}
		} else {
			select {
			case entry := <-new:
				fmt.Printf("From Indexer: %s\n", entry.ID)
				input <- entry
				continue
			case entry := <-remotes:
				fmt.Printf("From Database: %s\n", entry.ID)
				input <- entry
				continue
			default:
				fmt.Printf("No Signal\n")
				time.Sleep(15 * time.Second)
			}
		}

	}
}

// Generate the number of worker processes to optimize efficiency.
// Subtract 2 from the number of cores because of the main thread and the GetAll function.
func genNumWorkers() int {
	if runtime.NumCPU() > 2 {
		return runtime.NumCPU() - 1
	}
	return 1
}

func runWorker(keysets map[string]int, input <-chan database.FileKey, output chan<- database.FileKey, num int) {
	for key := range input {
		fmt.Printf("[Worker %d] Busy\n", num)
		threshold := keysets[key.KeySet]
		replications, err := ipfs.FindProvs(key.ID, threshold)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("File: %s is backed up %d time(s) and the threshold is %d.\n", key.ID, replications, threshold)

		// Determine an at risk file.
		// Node: if a file is hosted 0 times don't try to pin it.
		if replications < threshold && replications >= 1 {
			key, err = ReplicateAtRiskFile(key, keysets, output)
			if err != nil {
				log.Fatal(err)
			}
		}
		key.Replications = replications
		output <- key
		fmt.Printf("[Worker %d] Free\n", num)
	}
}
