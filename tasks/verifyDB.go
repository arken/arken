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
			err = keysets.IndexFull(db, location, new, output)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("[Finished Indexing & Updating Keysets]")
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
