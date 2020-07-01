package engine

import (
	"fmt"
	"path/filepath"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
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
		// Pin Lighthouse File to determine the size of the active cluster.
		fmt.Println("Pinning Lighthouse File...")
		err = ipfs.Pin(config.Keysets[set].LightHouseFileID)
		if err != nil {
			return err
		}

		keySet := filepath.Base(config.Keysets[set].URL)

		fmt.Printf("Calculating File Minimum Nodes Threshold for %s\n", keySet)
		threshold, err := CalcThreshold(config.Keysets[set].LightHouseFileID, config.Keysets[set].ReplicationFactor, 20)
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
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		for key := range input {
			err := ReplicateAtRiskFile(tx, key, threshold)
			if err != nil {
				return err
			}
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
