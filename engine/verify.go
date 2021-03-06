package engine

import (
	"log"
	"os"
	"path/filepath"

	"github.com/arken/arken/ipfs"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
)

// VerifyLocal verifies to ipfs that the local files are still
// locally pinned to the system.
func VerifyLocal() {
	// Copy Database because we can't guarentee this won't run as something is added.
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "verify.db")
	err := database.Copy(config.Global.Database.Path, copyName)
	if err != nil {
		log.Fatal(err)
	}

	// Open Copy Database
	read, err := database.Open(copyName)
	if err != nil {
		log.Fatal(err)
	}

	input := make(chan database.FileKey)
	// Get all will read db entries and put in queue for workers.
	go database.GetAll(read, "local", "", input)

	// Iterate through all local pins to verify pinned status.
	for entry := range input {
		err = ipfs.Pin(entry.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = read.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(copyName)
	if err != nil {
		log.Fatal(err)
	}
}
