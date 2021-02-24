package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/ipfs"
	"github.com/arken/arken/keysets"
)

func verifyDB(keySets []config.KeySet, new chan database.FileKey, output chan database.FileKey) {
	// Run Verify After Indexing has Finished
	for {
		config.Locks.IndexingSets.Lock()
		fmt.Println("\n[Verifying Internal DB against Keysets]")

		// Verify files in IPFS node match with Database before copying the database.
		nodeFiles := make(chan database.FileKey)
		go ipfs.LsPin(nodeFiles)
		for entry := range nodeFiles {
			output <- entry
		}

		startTime := time.Now()
		// Bring the modified date past the current check time for all files still within a keyset.
		for _, keySet := range keySets {
			err := keysets.IndexFull(keySet.Path, output)
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
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				continue
			}
			log.Fatal(err)
		}

		// Check each keyset for files that were not touched during the index
		files := make(chan database.FileKey)
		go database.GetAll(db, "local+remote", "", files)

		for entry := range files {
			if entry.Modified.Before(startTime) {
				entry.Status = "removed"
				output <- entry
				if config.Flags.Verbose {
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
			if config.Flags.Verbose {
				fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
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
		config.Locks.IndexingSets.Unlock()
		time.Sleep(7 * 24 * time.Hour)
	}
}
