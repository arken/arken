package engine

import (
	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

func makeSpace(bytes int64, keysets map[string]int, output chan<- database.FileKey) (removedBytes int64, err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return -1, err
	}
	// Set up database channels.
	input := make(chan database.FileKey)
	signal := make(chan bool)
	go database.GetStream(db, "local", "", input, signal)

	// Create Sum to track removed files.
	var sum int64

	// Iterate through entries returned through channel
	signal <- true
	for entry := range input {
		if entry.Name == "lighthouse" {
			signal <- true
			continue
		}
		replications, err := ipfs.FindProvs(entry.ID, keysets[entry.KeySet]+10)
		if err != nil {
			return -1, err
		}
		if replications > keysets[entry.KeySet] {
			// Unpin the file from storage.
			err := ipfs.Unpin(entry.ID)
			if err != nil {
				return -1, err
			}

			// Update value in database.
			entry.Status = "unpinned"
			output <- entry

			// Record updated removed sum.
			sum = sum + int64(entry.Size)
		}
		if sum >= bytes {
			signal <- false
		} else {
			signal <- true
		}
	}
	return sum, nil
}
