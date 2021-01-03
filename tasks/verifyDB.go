package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
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

		// Verify files in IPFS node match with Database before copying the database.
		nodeFiles := make(chan database.FileKey)
		go ipfs.LsPin(nodeFiles)
		for entry := range nodeFiles {
			output <- entry
		}
		// Wait for Writer to Finish Processing Pins
		for len(output) > 0 {
			time.Sleep(1 * time.Second)
		}

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
			fmt.Printf("Verifying: %s\n", filepath.Base(keySets[keySet].URL))

			location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keySets[keySet].URL))
			err = keysets.IndexFull(db, location, nil, output)
			if err != nil {
				log.Fatal(err)
			}

			lighthouse, err := database.Get(db, keySets[keySet].LightHouseFileID)
			if err != nil {
				log.Fatal(err)
			}
			output <- lighthouse

			remotes := make(chan database.FileKey)
			go database.GetAll(db, "local+remote", lighthouse.KeySet, remotes)

			for entry := range remotes {
				if lighthouse.Modified.After(entry.Modified) && entry.Name != "lighthouse" {
					entry.Status = "removed"
					output <- entry
					fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
				}
			}
		}
		orphans := make(chan database.FileKey)
		// Remove Orphaned Files without a Keyset.
		go database.GetAll(db, "local+remote", "_", orphans)
		for entry := range orphans {
			entry.Status = "removed"
			output <- entry
			fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
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
