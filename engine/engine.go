package engine

import (
	"fmt"
	"log"
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

	for set := range config.Keysets.Sets {
		keySet := filepath.Base(config.Keysets.Sets[set])

		fmt.Printf("Calculating File Minimum Nodes Threshold for %s\n", keySet)
		threshold, err := CalcThreshold("QmQBywyRvS3MJCP8jbV4Bsz8WMbRFsoux6EjsEwHhBDWqe", 2)
		if err != nil {
			return err
		}

		err = ScanHostReplications(db, keySet, threshold)
		if err != nil {
			return err
		}

		input := make(chan database.FileKey)
		go database.GetAll(db, "atrisk", keySet, input)
		for key := range input {
			err := ReplicateAtRiskFile(db, key, threshold)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}
