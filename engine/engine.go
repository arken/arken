package engine

import (
	"path/filepath"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/database"
)

// Rebalance manages balancing new and at risk files
// between nodes.
func Rebalance() (err error) {
	db, err := database.Open(config.Global.Database.Path)

	for set := range config.Keysets.Sets {
		err = ScanHostReplications(db, filepath.Base(config.Keysets.Sets[set]))
		if err != nil {
			return err
		}
	}

	return nil
}
