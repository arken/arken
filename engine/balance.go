package engine

import (
	"path/filepath"
	"strings"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

func makeSpace(bytes int64) (removedBytes int64, err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return -1, err
	}

	// Cache file replication thresholds for all keysets.
	thresholds := make(map[string]int)
	for _, keyset := range config.Keysets {
		name := strings.Split(filepath.Base(keyset.URL), ".")[0]
		hold, err := CalcThreshold(keyset.LightHouseFileID, keyset.ReplicationFactor, 20)
		if err != nil {
			return -1, err
		}
		thresholds[name] = hold
	}

	// Set up database channels.
	input := make(chan database.FileKey)
	signal := make(chan bool)
	go database.GetStream(db, "local", "", input, signal)

	// Create Sum to track removed files.
	var sum int64

	// Create Database Transaction
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// Iterate through entries returned through channel
	signal <- true
	for entry := range input {
		replications, err := ipfs.FindProvs(entry.ID, thresholds[entry.KeySet]+10)
		if err != nil {
			return -1, err
		}
		if replications > thresholds[entry.KeySet] {
			// Unpin the file from storage.
			err := ipfs.Unpin(entry.ID)
			if err != nil {
				return -1, err
			}

			// Update value in database.
			entry.Status = "remote"
			database.Update(tx, entry)

			// Record file operation/transaction
			database.TransactionCommit(tx, "removed", entry)

			// Record updated removed sum.
			sum = sum + int64(entry.Size)
		}
		if sum >= bytes {
			signal <- false
		} else {
			signal <- true
		}
	}
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return sum, nil
}
