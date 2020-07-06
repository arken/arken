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
	updateTime(db, input)
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

func updateTime(db *sql.DB, entry FileKey) (err error) {
	err = db.Ping()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`UPDATE keys SET
			modified = datetime('now')
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.ID)

	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
