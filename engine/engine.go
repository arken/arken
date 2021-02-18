package engine

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/ipfs"
)

// NetworkLimit is true if the node has hit it's download limit for the month.
var NetworkLimit bool
var keysets map[string]int

// Run manages balancing new and at risk files
// between nodes.
func Run(new, remotes, output chan database.FileKey) (err error) {
	keysets = make(map[string]int)
	input := make(chan database.FileKey, 10)

	for set := range config.Keysets {
		name := strings.Split(filepath.Base(config.Keysets[set].URL), ".")[0]
		keysets[name] = config.Keysets[set].Replications
	}

	// Determine the possible number of threads for the system's CPU
	workers := genNumWorkers()
	// Generate Worker Threads
	for i := 0; i < workers; i++ {
		go runWorker(keysets, input, output, i)
	}

	for {
		if NetworkLimit {
			select {
			case entry := <-new:
				output <- entry
				continue
			case entry := <-remotes:
				output <- entry
				continue
			}
		} else {
			select {
			case entry := <-new:
				input <- entry
				continue
			case entry := <-remotes:
				input <- entry
				continue
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
		threshold := keysets[key.KeySet]
		replications, err := ipfs.FindProvs(key.ID, threshold)
		if err != nil {
			log.Fatal(err)
		}
		if config.Flags.Verbose {
			fmt.Printf("File: %s is backed up %d time(s) and the threshold is %d.\n", key.ID, replications, threshold)
		}

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
	}
}
