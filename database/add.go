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
	}
	return nil
}

// Insert adds a Keyset file entry to the database.
func Insert(db *sql.DB, entry FileKey) {
	// Ping to check that database connection still exists.
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare(
		`INSERT INTO keys(
			id,
			name,
			size,
			status
		) VALUES(?,?,?,?);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		entry.ID,
		entry.Name,
		entry.Size,
		entry.Status)

	if err != nil {
		log.Fatal(err)
	}
}
