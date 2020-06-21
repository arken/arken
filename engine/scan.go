package engine

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/archivalists/arken/database"
	"github.com/archivalists/arken/ipfs"
)

// ScanHostReplications scans remote files from imported keysets and queries
// the ipfs network for the number of peers hosting that file.
func ScanHostReplications(db *sql.DB, keySet string, threshold int) (err error) {

	fmt.Printf("Calculated Threshold To Be: %d\n", threshold)

	input := make(chan database.FileKey)
	atRisk := make(chan database.FileKey)

	workers := genNumWorkers()

	// Create Sync WaitGroup to watch goroutines so that we don't close the channels early.
	var wg sync.WaitGroup
	wg.Add(workers)

	// Get all will read db entries and put in queue for workers.
	go database.GetAll(db, "remote", keySet, input)
	for i := 0; i < workers; i++ {
		go runWorker(&wg, threshold, input, atRisk)
	}

	// Create a go routine to wait till all workers are finished before closing channel.
	go func() {
		wg.Wait()
		close(atRisk)
	}()

	// Update all db entires that are out-of-date.
	for key := range atRisk {
		fmt.Println(key.ID)
		database.Update(db, key)
	}

	return nil
}

// Generate the number of worker processes to optimize efficiency.
// Subtract 2 from the number of cores because of the main thread and the GetAll function.
func genNumWorkers() int {
	if runtime.NumCPU() > 2 {
		return runtime.NumCPU() - 2
	}
	return 1
}

func runWorker(wg *sync.WaitGroup, threshold int, input <-chan database.FileKey, output chan<- database.FileKey) {
	for key := range input {
		replications, err := ipfs.FindProvs(key.ID, threshold)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("File: %s is backed up %d time(s) and the threshold is %d.\n", key.ID, replications, threshold)

		if replications < threshold {
			key.Status = "AtRisk"
			output <- key
		}
	}
	wg.Done()
}
