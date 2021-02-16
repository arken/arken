package engine

import (
	"fmt"
	"math/rand"

	"github.com/arken/arken/config"

	"github.com/arken/arken/database"
	"github.com/arken/arken/ipfs"
)

/*
ReplicateAtRiskFile will pin a file in danger of being lost to local storage.
This function will also run the El Farol Mathematics Problem to determine the
probability that this node should grab the file
*/
func ReplicateAtRiskFile(file database.FileKey, keysets map[string]int, write chan<- database.FileKey) (output database.FileKey, err error) {
	activationEnergy := float32(file.Replications) / float32(keysets[file.KeySet])
	prob := rand.Float32()

	if prob > activationEnergy {
		file.Size, err = ipfs.GetSize(file.ID)
		if err != nil {
			return file, nil
		}

		// If the file is bigger than the entire pool size then don't try to pin it.
		poolSize := config.ParseWellFormedPoolSize(config.Global.General.PoolSize)
		if uint64(file.Size) >= poolSize {
			return file, nil
		}

		repoSize, err := ipfs.GetRepoSize()
		if err != nil {
			if err.Error() == "context deadline exceeded" {
				return file, nil
			}
			return file, err
		}
		if int64(file.Size) > int64(poolSize)-int64(repoSize) {
			err := MakeSpace(int64(file.Size), write, false)
			if err != nil {
				if err.Error() == "could not make space" {
					return file, nil
				}
				return file, err
			}
		}

		fmt.Printf("Pinning to Local Storage: %s\n", file.ID)
		err = ipfs.Pin(file.ID)
		if err != nil {
			return file, err
		}
		ipfs.AdjustRepoSize(int64(file.Size))
		fmt.Printf("Pinned to Local Storage: %s\n", file.ID)

		file.Status = "local"
		return file, nil
	}
	return file, nil
}
