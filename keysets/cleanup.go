package keysets

import (
	"fmt"
	"path/filepath"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

func garbageCollect(keySet config.KeySet) (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	// Wait to close the database until all files have been indexed.
	defer db.Close()

	lighthouse, err := database.Get(db, keySet.LightHouseFileID)
	if err != nil {
		return err
	}

	input := make(chan database.FileKey)

	go database.GetAll(db, "local+remote", filepath.Base(keySet.URL), input)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	for key := range input {
		if key.Modified.Before(lighthouse.Modified) {
			database.Delete(tx, key.ID)
			if key.Status == "local" {
				err = ipfs.Unpin(key.ID)
				if err != nil {
					return err
				}
			}
			fmt.Printf("Removed: %s\n", key.ID)
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
