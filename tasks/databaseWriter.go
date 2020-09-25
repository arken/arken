package tasks

import (
	"database/sql"
	"log"
	"time"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

func databaseWriter(input chan database.FileKey, settings chan string) {
	var (
		timeout int
		db      *sql.DB
		err     error
	)

	for {
		select {
		// If a commit is recieved write it's update to the database.
		case commit := <-settings:
			// If the database has been closed due to a timeout of inactivity.
			// reopen the connection to the database.
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
			continue

		// If an entry is recieved triage it and write it to the database.
		case entry := <-input:
			// Reset timeout on signal recieved.
			timeout = 0
			if db == nil {
				db, err = database.Open(config.Global.Database.Path)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Check for previous entry in database.
			prev, err := database.Get(db, entry.ID)
			if err != nil && err.Error() != "entry not found" {
				log.Fatal(err)
			}
			switch {
			case prev.Status == "":
				switch {
				case entry.Name == "lighthouse":
					database.Add(db, entry)
					entry.Status = "local"
					database.Update(db, entry)
				case entry.Status == "local":
					database.Insert(db, entry)
					database.TransactionCommit(db, "added", entry)
					continue

				case entry.Status == "removed":
					continue

				// Cover "added", "remote", and "unpinned" statuses.
				default:
					entry.Status = "remote"
					database.Insert(db, entry)
					continue
				}
			case prev.Status == "local":
				switch {
				case entry.Status == "removed":
					ipfs.Unpin(entry.ID)
					database.Update(db, entry)
					continue

				case entry.Status == "unpinned":
					entry.Status = "remote"
					database.Update(db, entry)
					continue

				// Cover "added", "local", "remote"
				default:
					database.UpdateTime(db, entry)
					continue
				}
			case prev.Status == "remote":
				switch {
				case entry.Status == "local":
					database.Update(db, entry)
					database.TransactionCommit(db, "added", entry)
					continue

				case entry.Status == "removed":
					database.Update(db, entry)
					continue

				// Cover "added", "remote", "unpinned"
				default:
					database.UpdateTime(db, entry)
					continue
				}
			case prev.Status == "removed":
				switch {
				case entry.Status == "added":
					entry.Status = "remote"
					database.Update(db, entry)
					continue

				case entry.Status == "local":
					ipfs.Unpin(entry.ID)
					database.TransactionCommit(db, "removed", entry)
					database.Delete(db, entry.ID)
					continue

				case entry.Status == "removed":
					database.UpdateTime(db, entry)
					continue

				// Cover "remote", "unpinned"
				default:
					database.Delete(db, entry.ID)
					continue
				}
			}

		default:
			if timeout > 30 && db != nil {
				db.Close()
				db = nil
			} else if timeout > 30 {
				time.Sleep(30 * time.Second)
			} else if timeout > 15 {
				time.Sleep(5 * time.Second)
				timeout++
			} else {
				timeout++
			}
		}
	}
}
