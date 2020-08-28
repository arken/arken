package engine

import (
	"io"
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
	dbFile, err := os.Open(config.Global.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	// Copy Database because we can't guarentee this won't run as something is added.
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "verifyF.db")
	copyFile, err := os.Create(copyName)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(copyName)

	_, err = io.Copy(copyFile, dbFile)
	if err != nil {
		log.Fatal(err)
	}
	copyFile.Close()
	dbFile.Close()

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
		_, err = ipfs.GetSize(entry.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
