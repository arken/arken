package engine

import (
	"fmt"
	"log"
)

func (n *Node) VerifyDB() {
	fmt.Printf("Starting Database Verification...\n")
	// Check for missmatches between database and manifest
	files, err := n.Manifest.Index(n.DB, true)
	if err != nil {
		log.Println(err)
	}

	// Boot adder subsystem
	toAdder, err := n.FileAdder()
	if err != nil {
		log.Println(err)
	}

	for file := range files {
		// If file has been added to the manifest, add it
		// to the database and check number of times replicated
		// to determine if file should be replicated locally.
		if file.Status == "add" {
			file.Status = "remote"

			// Add file to database.
			err = n.DB.Add(file)
			if err != nil {
				log.Println(err)
				continue
			}

			// Send file to adder subsystem
			toAdder <- file
		}

		// If file has been deleted from the manifest, remove
		// it from the database and unpin it from the embedded
		// IPFS node if necessary.
		if file.Status == "remove" {
			result, err := n.DB.Remove(file.ID)
			if err != nil {
				log.Println(err)
				continue
			}

			if result.Status == "local" {
				err = n.Node.Unpin(result.ID)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}

	// Close adder subsystem
	close(toAdder)

	// Clean up IPFS store after syncing manifest.
	err = n.Node.GC()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Database Verification Complete\n")
}
