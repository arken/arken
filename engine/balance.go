package engine

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

var (
	output chan []database.FileKey
	signal chan spaceRequest
)

type spaceRequest struct {
	space int64
	force bool
}

func init() {
	signal = make(chan spaceRequest)
	output = make(chan []database.FileKey)

	go makeSpaceDaemon()
}

func makeSpaceDaemon() {
	var (
		running  bool
		timeout  int
		db       *sql.DB
		err      error
		ping     chan bool
		input    chan database.FileKey
		copyName string
		lastSum  int64
	)

	for {
		select {
		case request := <-signal:
			if db == nil {
				copyName = filepath.Join(filepath.Dir(config.Global.Database.Path), "collect.db")
				err = database.Copy(config.Global.Database.Path, copyName)
				if err != nil {
					log.Fatal(err)
				}
				db, err = database.Open(copyName)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				if lastSum == 0 {
					output <- []database.FileKey{}
					timeout = 0
					continue
				}
			}
			if !running {
				// Set up database channels.
				input = make(chan database.FileKey)
				ping = make(chan bool)
				go database.GetStream(db, "local", "", input, ping)
				running = true
			}
			// Send request to GetStream
			ping <- true
			// Create Sum to track removed files.
			var sum int64
			sum = 0
			var response []database.FileKey

			for entry := range input {
				if entry.Name == "lighthouse" {
					ping <- true
					continue
				}
				replications, err := ipfs.FindProvs(entry.ID, keysets[entry.KeySet]+10)
				if err != nil {
					log.Fatal(err)
				}
				if request.force || replications > keysets[entry.KeySet] {
					response = append(response, entry)

					// Record updated removed sum.
					sum = sum + int64(entry.Size)
				}
				if sum >= request.space {
					break
				} else {
					ping <- true
				}
			}
			if sum >= request.space {
				output <- response
			} else {
				running = false
				output <- []database.FileKey{}
			}
			lastSum = sum
			timeout = 0
			continue

		default:
			if timeout > 30 && db != nil {
				ping <- false
				db.Close()
				db = nil
				lastSum = -1
				running = false
				os.Remove(copyName)
			} else {
				timeout++
			}
			time.Sleep(15 * time.Second)
		}
	}
}

// MakeSpace unpins X bytes worth of files from the node.
// In the case of adding new files this is done non-forcefully to only remove
// well backed up files if possible.
// In the case of forcing the node back to the pool size this is done forcefully
// removing all files nessisary to meet the request.
func MakeSpace(bytes int64, filesOut chan<- database.FileKey, force bool) (err error) {
	// Make request to backend daemon.
	signal <- spaceRequest{space: bytes, force: force}

	// Wait for response.
	response := <-output

	if len(response) > 0 {
		for _, file := range response {
			// Unpin the file from storage.
			err := ipfs.Unpin(file.ID)
			if err != nil {
				return err
			}
			file.Status = "unpinned"
			filesOut <- file
		}
		return nil
	}
	return errors.New("could not make space")
}
