package engine

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/arkenproject/arken/config"

	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

/*
ReplicateAtRiskFile will pin a file in danger of being lost to local storage.
This function will also run the El Farol Mathematics Problem to determine the
probability that this node should grab the file
*/
func ReplicateAtRiskFile(tx *sql.Tx, file database.FileKey, threshold int) (err error) {
	replications, err := ipfs.FindProvs(file.ID, threshold)
	if err != nil {
		return err
	}
	activationEnergy := float32(replications) / float32(threshold)
	prob := rand.Float32()

	if prob > activationEnergy {
		file.Size, err = ipfs.GetSize(file.ID)
		if err != nil {
			return err
		}

		if uint64(file.Size) >= config.Disk.GetPoolSizeBytes() {
			return nil
		}

		// To Do: Add the logic here for removing well backed up files in favor of at risk files.
		if uint64(file.Size) >= config.Disk.GetAvailableBytes() {
			return nil
		}

		fmt.Printf("Pinning to Local Storage: %s\n", file.ID)
		err = ipfs.Pin(file.ID)
		if err != nil {
			return err
		}
		database.TransactionCommit(tx, "added", file)
		file.Status = "local"
	} else {
		file.Status = "remote"
	}
	database.Update(tx, file)
	return nil
}
