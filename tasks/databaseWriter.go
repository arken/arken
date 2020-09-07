package tasks

import (
	"database/sql"
	"log"

	"github.com/arkenproject/arken/database"
)

func databaseWriter(db *sql.DB, input chan database.FileKey) {
	for entry := range input {
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

		if entry.Status == "remote" {
			database.UpdateTime(db, entry)
			continue
		}

		if entry.Status == "removed" {
			database.Delete(db, entry.ID)
			if prev.Status == "local" {
				database.TransactionCommit(db, "removed", entry)
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

	}
}
