package tasks

import (
	"database/sql"
	"log"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

func databaseWriter(input chan database.FileKey, settings chan string) {
	var (
		timeout int
		db      *sql.DB
		err     error
	)

	for {
		select {
		case commit := <-settings:
			if db == nil {
				db, err = database.Open(config.Global.Database.Path)
				if err != nil {
					log.Fatal(err)
				}
			}
			err := database.SetCommit(db, commit)
			if err != nil {
				log.Fatal(err)
			}
			timeout = 0

		case entry := <-input:
			if db == nil {
				db, err = database.Open(config.Global.Database.Path)
				if err != nil {
					log.Fatal(err)
				}
			}
			timeout = 0
			// Test if the File is in the database.
			prev, err := database.Get(db, entry.ID)
			if err != nil && err.Error() == "entry not found" {
				// If the entry is not found is should be new.
				if entry.Status != "removed" {
					database.Add(db, entry)
					if entry.Status == "local" {
						database.TransactionCommit(db, "added", entry)
					}
					continue
				} else {
					// This would be if the file is marked for
					// deletion without being in the database somehow.
					database.Delete(db, entry.ID)
					continue
				}
			}
			if err != nil {
				log.Fatal(err)
			}

			if prev.Status == "removed" {
				if entry.Status == "local" {
					err = ipfs.Unpin(entry.ID)
					if err != nil {
						log.Fatal(err)
					}
					database.TransactionCommit(db, "removed", entry)
				}
				database.Delete(db, entry.ID)
			}

			if entry.Status == "remote" {
				database.UpdateTime(db, entry)
				continue
			}

			if entry.Status == "removed" {
				if prev.Status == "local" {
					err = ipfs.Unpin(entry.ID)
					if err != nil {
						log.Fatal(err)
					}
					database.Delete(db, entry.ID)
					database.TransactionCommit(db, "removed", entry)
				} else {
					database.Update(db, entry)
				}
				continue
			}

			if entry.Status == "unpinned" {
				entry.Status = "remote"
				database.Update(db, entry)
				database.TransactionCommit(db, "removed", entry)
				continue
			}

			if entry.Status == "local" {
				database.Update(db, entry)
				database.TransactionCommit(db, "added", entry)
			}

		default:
			if timeout > 30 && db != nil {
				db.Close()
				db = nil
			} else {
				timeout++
			}
			time.Sleep(15 * time.Second)
		}
	}
}
