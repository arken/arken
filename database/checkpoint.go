package database

import (
	"database/sql"
	"errors"
	"log"
)

// SetCommit sets the hash provided as the checkpointed hash in the database.
func SetCommit(db *sql.DB, hash string) (err error) {
	_, err = GetCommit(db)
	if err != nil {
		if err.Error() == "entry not found" {
			insertCommit(db, hash)
		} else {
			return err
		}
	} else {
		err = updateCommit(db, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert adds a Keyset file entry to the database.
func insertCommit(db *sql.DB, hash string) {
	commit := "commit"
	stmt, err := db.Prepare(
		`INSERT INTO commitCheckpoint(
			name,
			hash
		) VALUES(?,?);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(commit, hash)

	if err != nil {
		log.Fatal(err)
	}
}

// GetCommit searches for and returns a the coorisponding entry from the
// database if the entry exists.
func GetCommit(db *sql.DB) (result string, err error) {
	commit := "commit"
	row, err := db.Query("SELECT hash FROM commitCheckpoint WHERE name = ?", commit)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		return result, errors.New("entry not found")
	}
	err = row.Scan(
		&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}

// updateCommit updates the checkpoint transaction of the entry in the DB
func updateCommit(db *sql.DB, hash string) (err error) {
	commit := "commit"
	stmt, err := db.Prepare(
		`UPDATE commitCheckpoint SET
			hash = ?
			WHERE name = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		hash,
		commit)

	return err
}
