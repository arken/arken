package engine

import (
	"fmt"
	"log"
	"os"

	"github.com/arken/arken/manifest"
	"github.com/go-git/go-git/v5"
)

func (n *Node) SyncManifest() {
	fmt.Printf("Starting Manifest Import...\n")

	// Pull changes from upstream manifest
	err := n.Manifest.Pull()
	if err != nil {
		// Check if a non-fast forward update occurred while pulling.
		if err == git.ErrNonFastForwardUpdate {
			// Remove manifest and re-clone.
			err = os.RemoveAll(n.Cfg.Manifest.Path)
			if err != nil {
				log.Println(err)
			}

			// Re-initialize the manifest
			n.Manifest, err = manifest.Init(
				n.Cfg.Manifest.Path,
				n.Cfg.Manifest.URL,
			)
		}
		if err != nil {
			log.Println(err)
		}
	}

	// Update manifest settings
	err = n.Manifest.Decode()
	if err != nil {
		log.Println(err)
	}

	// Index changes from manifest
	files, err := n.Manifest.Index(n.DB, false)
	if err != nil {
		log.Println(err)
	}

	// Check if there are any new files indexed
	// from the manifest.
	if file, ok := <-files; ok {
		// Put first file back in queue
		go func() { files <- file }()

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
	}

	fmt.Printf("Successfully imported & updated manifest\n")
}
