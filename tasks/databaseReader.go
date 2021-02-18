package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
)

func databaseReader(remote chan database.FileKey, output chan database.FileKey) {
	for {
		fmt.Println("\n[Starting Rebalancing]")
		copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "checkRemotes.db")
		database.Copy(config.Global.Database.Path, copyName)

		db, err := database.Open(copyName)
		if err != nil {
			log.Fatal(err)
		}

		// Grab all remote files and pass to engine to check for at risk files.
		relay := make(chan database.FileKey)
		go database.GetAll(db, "remote+added", "", relay)

		for entry := range relay {
			remote <- entry
		}

		deleted := make(chan database.FileKey)
		go database.GetAll(db, "removed", "", deleted)

		for entry := range deleted {
			entry.Status = "remote"
			output <- entry
		}

		err = os.Remove(copyName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n[Finished Data Rebalance]")

		time.Sleep(1 * time.Hour)
	}
}
