package engine

import (
	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/database"
)

// Rebalance manages balancing new and at risk files
// between nodes.
func Rebalance() (err error) {
	db, err := database.Open(config.Global.Database.Path)

	err = ScanHostReplications(db)
	if err != nil {
		return err
	}

	return nil
}
