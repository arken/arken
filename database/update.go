package database

import (
	"database/sql"
	"log"
	"time"
)

// Update changes a file's status in the database.
func Update(db *sql.DB, key FileKey) {
	stmt, err := db.Prepare(
		`UPDATE keys SET
			Status = ?,
			Replications = ?,
			Size = ?,
			Modified = ?
			WHERE id = ?;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		key.Status,
		key.Replications,
		key.Size,
		time.Now(),
		key.ID)

	if err != nil {
		log.Fatal(err)
	}
}
