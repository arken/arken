package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/keysets"

	"github.com/arkenproject/arken/database"
)

func loadSets(keySets []config.KeySet, new chan database.FileKey, output chan database.FileKey, settings chan string) {
	// Run LoadSets every hour.
	for {
		for config.Flags.IndexingSets {
			time.Sleep(15 * time.Minute)
		}
		fmt.Println("\n[Indexing & Updating Keysets]")
		config.Flags.IndexingSets = true
		err := keysets.LoadSets(keySets)
		if err != nil {
			log.Println(err)
		} else {
			for keySet := range keySets {
				location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keySets[keySet].URL))
				lighthouse, err := keysets.ConfigLighthouse(keySets[keySet].LightHouseFileID, keySets[keySet].URL)
				if err != nil {
					log.Fatal(err)
				}
				output <- lighthouse

				fmt.Printf("Indexing: %s\n", filepath.Base(keySets[keySet].URL))
				err = keysets.Index(location, new, output, settings)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Println("[Finished Indexing & Updating Keysets]")
		config.Flags.IndexingSets = false
		time.Sleep(1 * time.Hour)
	}
}
