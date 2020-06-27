package engine

import (
	"fmt"
	"path/filepath"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/database"
)

// Rebalance manages balancing new and at risk files
// between nodes.
func Rebalance() (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	defer db.Close()

	for set := range config.Keysets {
		keySet := filepath.Base(config.Keysets[set].URL)

		fmt.Printf("Calculating File Minimum Nodes Threshold for %s\n", keySet)
		threshold, err := CalcThreshold(config.Keysets[set].LightHouseFileID, config.Keysets[set].ReplicationFactor, 5)
		if err != nil {
			return err
		}

		fmt.Printf("Calculated Threshold To Be: %d\n", threshold)

		err = ScanHostReplications(db, keySet, threshold)
		if err != nil {
			return err
		}

		input := make(chan database.FileKey)
		go database.GetAll(db, "atrisk", keySet, input)
		for key := range input {
			err := ReplicateAtRiskFile(db, key, threshold)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
