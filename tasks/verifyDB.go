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
		time.Sleep(30 * time.Second)
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

		// Bring the modified date past the lighthouse file for all files still within a keyset.
		for _, keySet := range keySets {
			location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keySet.URL))
			lighthouse, err := keysets.ConfigLighthouse(keySet.LightHouseFileID, keySet.URL)
			if err != nil {
				log.Fatal(err)
			}
			output <- lighthouse

			err = keysets.IndexFull(location, output)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Wait for Writer to Finish Processing Pins
		for len(output) > 0 {
			time.Sleep(1 * time.Second)
		}

		// Copy the Database with the updated modified times.
		copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "verify.db")
		err := database.Copy(config.Global.Database.Path, copyName)
		if err != nil {
			log.Fatal(err)
		}
		db, err := database.Open(copyName)
		if err != nil {
			log.Fatal(err)
		}

		// Check each keyset for files that were not touched during the index
		for _, keySet := range keySets {
			lighthouse, err := database.Get(db, keySet.LightHouseFileID)
			if err != nil {
				log.Fatal(err)
			}

			files := make(chan database.FileKey)
			go database.GetAll(db, "local+remote", lighthouse.KeySet, files)

			for entry := range files {
				if entry.Modified.Before(lighthouse.Modified) {
					entry.Status = "removed"
					output <- entry
					fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
				}
			}
		}

		// Remove files found by IPFS with no owning Keyset.
		orphans := make(chan database.FileKey)
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
