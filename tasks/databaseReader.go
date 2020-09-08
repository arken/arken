package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
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

		relay := make(chan database.FileKey, 10)
		go database.GetAll(db, "remote", "", relay)

		for entry := range relay {
			remote <- entry
		}

		deleted := make(chan database.FileKey, 10)
		go database.GetAll(db, "removed", "", deleted)

		for entry := range deleted {
			output <- entry
		}

		os.Remove(copyName)
		fmt.Println("\n[Finished Data Rebalance]")

		time.Sleep(1 * time.Hour)
	}
}
