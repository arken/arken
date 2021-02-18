package tasks

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/arken/arken/config"
	"github.com/arken/arken/keysets"

	"github.com/arken/arken/database"
)

func loadSets(keySets []config.KeySet, added chan<- database.FileKey, output chan<- database.FileKey, settings chan string) {
	// Run LoadSets every hour.
	for {
		config.Locks.IndexingSets.Lock()
		fmt.Println("\n[Indexing & Updating Keysets]")
		err := keysets.LoadSets(keySets)
		if err != nil {
			log.Println(err)
		} else {
			for keySet := range keySets {
				location := keySets[keySet].Path
				if config.Flags.Verbose {
					fmt.Printf("Indexing: %s\n", filepath.Base(keySets[keySet].URL))
				}
				err = keysets.Index(location, added, output, settings)

				if err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Println("[Finished Indexing & Updating Keysets]")
		config.Locks.IndexingSets.Unlock()
		time.Sleep(1 * time.Hour)
	}
}
