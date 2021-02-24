package tasks

import (
	"database/sql"
	"log"
	"time"

	"github.com/arken/arken/ipfs"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
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
			case prev.Status == "" || prev.Status == "added":
				switch {
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
					err := ipfs.Unpin(entry.ID)
					if err != nil {
						database.Update(db, prev)
						continue
					}
					ipfs.AdjustRepoSize(0 - int64(prev.Size))
					database.TransactionCommit(db, "removed", entry)
					database.Update(db, entry)
					continue

				case entry.Status == "unpinned":
					entry.Status = "remote"
					database.Update(db, entry)
					continue

				case entry.Status == "added":
					prev.Name = entry.Name
					database.Update(db, prev)

				// Cover "local", "remote"
				default:
					database.Update(db, prev)
					continue
				}
			case prev.Status == "remote":
				switch {
				case entry.Status == "local":
					entry.Name = prev.Name
					database.Update(db, entry)
					database.TransactionCommit(db, "added", entry)
					continue

				case entry.Status == "removed":
					database.Update(db, entry)
					continue

				// Cover "added", "remote", "unpinned"
				default:
					database.Update(db, prev)
					continue
				}
			case prev.Status == "removed":
				switch {
				case entry.Status == "added":
					entry.Status = "remote"
					database.Update(db, entry)
					continue

				case entry.Status == "local":
					err := ipfs.Unpin(entry.ID)
					if err != nil {
						database.Update(db, entry)
						continue
					}
					ipfs.AdjustRepoSize(0 - int64(entry.Size))
					database.TransactionCommit(db, "removed", entry)
					database.Delete(db, entry.ID)
					continue

				case entry.Status == "removed":
					database.Update(db, prev)
					continue

				// Cover "remote", "unpinned"
				default:
					database.Delete(db, entry.ID)
					continue
				}
			}

		default:
			if timeout > 75 && db != nil {
				db.Close()
				db = nil
				ipfs.GC()
			} else if timeout > 75 {
				time.Sleep(30 * time.Second)
			} else if timeout > 60 {
				time.Sleep(5 * time.Second)
				timeout++
			} else if timeout > 30 {
				time.Sleep(200 * time.Millisecond)
				timeout++
			} else {
				time.Sleep(5 * time.Millisecond)
				timeout++
			}
		}
	}
}
