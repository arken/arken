package keysets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

// GarbageCollect looks for files that should be removed from
// the internal database.
func GarbageCollect(keySet config.KeySet, output chan database.FileKey) (err error) {
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "gc.db")
	err = database.Copy(config.Global.Database.Path, copyName)

	db, err := database.Open(copyName)
	if err != nil {
		os.Remove(copyName)
		return err
	}
	// Wait to close the database until all files have been indexed.
	defer db.Close()

	lighthouse, err := database.Get(db, keySet.LightHouseFileID)
	for err != nil {
		if strings.HasPrefix(err.Error(), "entry not found") {
			time.Sleep(5 * time.Second)
			err = database.Copy(config.Global.Database.Path, copyName)
			if err != nil {
				os.Remove(copyName)
				return err
			}
			lighthouse, err = database.Get(db, keySet.LightHouseFileID)
		} else {
			os.Remove(copyName)
			return err
		}
	}

	input := make(chan database.FileKey)

	go database.GetAll(db, "local+remote", filepath.Base(keySet.URL), input)
	for key := range input {
		if key.Modified.Before(lighthouse.Modified) {
			key.Status = "removed"
			if key.Status == "local" {
				err = ipfs.Unpin(key.ID)
				if err != nil {
					os.Remove(copyName)
					return err
				}
				key.Status = "removed"
			}
			output <- key
			fmt.Printf("Removed: %s\n", key.ID)
		}
	}
	os.Remove(copyName)
	return nil
}
