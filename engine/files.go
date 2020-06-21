package engine

import (
	"math/rand"

	"github.com/archivalists/arken/ipfs"
)

/*
ReplicateAtRiskFile will pin a file in danger of being lost to local storage.
This function will also run the El Farol Mathematics Problem to determine the
probability that this node should grab the file
*/
func ReplicateAtRiskFile(fileID string, threshold int) (err error) {
	replications, err := ipfs.FindProvs(fileID, threshold)
	if err != nil {
		return err
	}
	activationEnergy := float64(replications) / float64(threshold)
	prob := rand.Float64()

	if prob > activationEnergy {
		err = ipfs.Pin(fileID)
		if err != nil {
			return err
		}
	}

	return nil
}
