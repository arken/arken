package engine

import (
	"log"
	"os"
	"path/filepath"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

// VerifyLocal verifies to ipfs that the local files are still
// locally pinned to the system.
func VerifyLocal() {
	// Copy Database because we can't guarentee this won't run as something is added.
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "verify.db")
	database.Copy(config.Global.Database.Path, copyName)

	// Open Copy Database
	read, err := database.Open(copyName)
	if err != nil {
		log.Fatal(err)
	}
	defer read.Close()

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
	os.Remove(copyName)
}
