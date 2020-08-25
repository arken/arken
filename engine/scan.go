package engine

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

// ScanReplications scans remote files from imported keysets and queries
// the ipfs network for the number of peers hosting that file.
func ScanReplications(read *sql.DB, write *sql.DB, keySet string, threshold int) (err error) {
	input := make(chan database.FileKey)
	atRisk := make(chan database.FileKey)

	workers := genNumWorkers()

	// Create Sync WaitGroup to watch goroutines so that we don't close the channels early.
	var wg sync.WaitGroup
	wg.Add(workers)

	// Get all will read db entries and put in queue for workers.
	go database.GetAll(read, "remote", keySet, input)
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
		if key.Status == "local" {
			tx, err := write.Begin()
			if err != nil {
				return err
			}

			database.Update(tx, key)
			database.TransactionCommit(tx, "added", key)

			err = tx.Commit()
			if err != nil {
				return err
			}

		}
	}
	return nil
}

// Generate the number of worker processes to optimize efficiency.
// Subtract 2 from the number of cores because of the main thread and the GetAll function.
func genNumWorkers() int {
	if runtime.NumCPU() > 2 {
		return runtime.NumCPU() - 1
	}
	return 1
}

func runWorker(wg *sync.WaitGroup, threshold int, input <-chan database.FileKey, output chan<- database.FileKey) {
	for key := range input {
		replications, err := ipfs.FindProvs(key.ID, threshold)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("File: %s is backed up %d time(s) and the threshold is %d.\n", key.ID, replications, threshold)

		// Determine an at risk file.
		// Node: if a file is hosted 0 times don't try to pin it.
		if replications < threshold && replications >= 1 {
			key, err = ReplicateAtRiskFile(key, threshold)
			if err != nil {
				log.Fatal(err)
			}
		}
		key.Replications = replications
		output <- key
	}
	wg.Done()
}
