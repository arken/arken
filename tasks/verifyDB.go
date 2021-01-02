package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/keysets"
)

func verifyDB(keySets []config.KeySet, new chan database.FileKey, output chan database.FileKey) {
	// Run Verify After Indexing has Finished
	for {
		time.Sleep(1 * time.Minute)
		for config.Flags.IndexingSets {
			time.Sleep(1 * time.Hour)
		}
		fmt.Println("\n[Verifying Internal DB against Keysets]")
		config.Flags.IndexingSets = true
		copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "verifyKeyset.db")
		err := database.Copy(config.Global.Database.Path, copyName)
		if err != nil {
			log.Fatal(err)
		}

		db, err := database.Open(copyName)
		if err != nil {
			log.Fatal(err)
		}

		for keySet := range keySets {
			location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keySets[keySet].URL))
			lighthouse, err := keysets.ConfigLighthouse(keySets[keySet].LightHouseFileID, keySets[keySet].URL)
			if err != nil {
				log.Fatal(err)
			}
			output <- lighthouse

			fmt.Printf("Verifying: %s\n", filepath.Base(keySets[keySet].URL))
			err = keysets.IndexFull(db, location, nil, output)
			if err != nil {
				log.Fatal(err)
			}

			remotes := make(chan database.FileKey)
			go database.GetAll(db, "local+remote", lighthouse.KeySet, remotes)

			for entry := range remotes {
				if entry.Modified.After(lighthouse.Modified) {
					entry.Status = "removed"
					output <- entry
					fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
				}
			}
		}

		fmt.Println("[Finished Verifying Keysets]")
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(copyName)
		if err != nil {
			log.Fatal(err)
		}
		config.Flags.IndexingSets = false
		time.Sleep(7 * 24 * time.Hour)
	}
}
