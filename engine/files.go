package engine

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

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
	time.Sleep(1 * time.Second)

	activationEnergy := float32(file.Replications) / float32(threshold)
	prob := rand.Float32()

	if prob > activationEnergy {
		file.Size, err = ipfs.GetSize(file.ID)
		if err != nil {
			return err
		}

		// If the file is bigger than the entire pool size then don't try to pin it.
		poolSize := config.ParseWellFormedPoolSize(config.Global.General.PoolSize)
		if err != nil {
			return err
		}
		if uint64(file.Size) >= poolSize {
			return nil
		}

		// To Do: Add the logic here for removing "well backed" up files in favor of "at risk" files.
		repoSize, err := ipfs.GetRepoSize()
		if err != nil {
			return err
		}
		if uint64(file.Size) >= poolSize-repoSize {
			bytes, err := makeSpace(int64(file.Size))
			if err != nil {
				return err
			}
			if bytes < int64(file.Size) {
				return nil
			}
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
