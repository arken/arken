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
func ScanHostReplications(db *sql.DB) (err error) {
	input := make(chan database.FileKey)
	atRisk := make(chan database.FileKey)

	workers := genNumWorkers()

	// Create Sync WaitGroup to watch goroutines so that we don't close the channels early.
	var wg sync.WaitGroup
	wg.Add(workers)

	// Get all will read db entries and put in queue for workers.
	go database.GetAll(db, "remote", input)
	for i := 0; i < workers; i++ {
		go runWorker(&wg, input, atRisk)
	}

	// Create a go routine to wait till all workers are finished before closing channel.
	go func() {
		wg.Wait()
		close(atRisk)
	}()

	// Update all db entires that are out-of-date.
	for key := range atRisk {
		fmt.Println(key.ID)
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

func runWorker(wg *sync.WaitGroup, input <-chan database.FileKey, output chan<- database.FileKey) {
	for key := range input {
		replications, err := ipfs.FindProvs(key.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("File: %s is backed up at least %d time(s).\n", key.ID, replications)

		if replications < ipfs.AtRiskThreshhold {
			key.Status = "AtRisk"
			output <- key
		}
	}
	wg.Done()
}
