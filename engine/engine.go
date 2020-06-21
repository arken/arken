package engine

import (
	"log"
	"path/filepath"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/database"
)

// Rebalance manages balancing new and at risk files
// between nodes.
func Rebalance() (err error) {
	db, err := database.Open(config.Global.Database.Path)

	for set := range config.Keysets.Sets {
		keySet := filepath.Base(config.Keysets.Sets[set])

		threshold, err := CalcThreshold("QmZtmD2qt6fJot32nabSP3CUjicnypEBz7bHVDhPQt9aAy", 2)
		if err != nil {
			return err
		}

		err = ScanHostReplications(db, keySet, threshold)
		if err != nil {
			return err
		}

		input := make(chan database.FileKey)
		go database.GetAll(db, "AtRisk", keySet, input)
		for key := range input {
			err := ReplicateAtRiskFile(key.ID, threshold)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}
