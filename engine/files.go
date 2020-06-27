package engine

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

/*
ReplicateAtRiskFile will pin a file in danger of being lost to local storage.
This function will also run the El Farol Mathematics Problem to determine the
probability that this node should grab the file
*/
func ReplicateAtRiskFile(db *sql.DB, file database.FileKey, threshold int) (err error) {
	replications, err := ipfs.FindProvs(file.ID, threshold)
	if err != nil {
		return err
	}
	activationEnergy := float32(replications) / float32(threshold)
	prob := rand.Float32()

	if prob > activationEnergy {
		fmt.Printf("Pinning to Local Storage: %s\n", file.ID)
		err = ipfs.Pin(file.ID)
		if err != nil {
			return err
		}
		file.Status = "local"
	} else {
		file.Status = "remote"
	}
	database.Update(db, file)

	return nil
}
