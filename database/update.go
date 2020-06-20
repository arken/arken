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
	stmt, err := db.Prepare(
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
