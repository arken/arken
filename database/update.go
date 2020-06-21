package database

import (
	"database/sql"
	"log"
)

// Update changes a file's status in the database.
func Update(db *sql.DB, key FileKey) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()
	stmt, err := tx.Prepare(
		`UPDATE keys SET
			Status = ?
			WHERE id = ?;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		key.Status,
		key.ID)

	if err != nil {
		log.Fatal(err)
	}
}
