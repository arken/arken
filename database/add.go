package database

import (
	"database/sql"
	"log"
)

// Add inserts an entry into the database if it doesn't exist already.
func Add(db *sql.DB, input FileKey) (err error) {
	_, err = Get(db, input.ID)
	if err != nil {
		if err.Error() == "entry not found" {
			Insert(db, input)
		} else {
			return err
		}
	} else {
		err = UpdateTime(db, input)
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert adds a Keyset file entry to the database.
func Insert(db *sql.DB, entry FileKey) {
	stmt, err := db.Prepare(
		`INSERT INTO keys(
			id,
			name,
			size,
			status,
			keyset,
			modified,
			replications
		) VALUES(?,?,?,?,?,datetime('now'),?);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		entry.ID,
		entry.Name,
		entry.Size,
		entry.Status,
		entry.KeySet,
		entry.Replications)

	if err != nil {
		log.Fatal(err)
	}
}

// UpdateTime sets the last seen time of a file to now in the database.
func UpdateTime(db *sql.DB, entry FileKey) (err error) {
	stmt, err := db.Prepare(
		`UPDATE keys SET
			modified = datetime('now')
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.ID)

	return err
}
