package database

import (
	"database/sql"
	"log"
)

// Update changes a file's status in the database.
func Update(tx *sql.Tx, key FileKey) {
	stmt, err := tx.Prepare(
		`UPDATE keys SET
			Status = ?,
			Replications = ?,
			Size = ?
			WHERE id = ?;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		key.Status,
		key.Replications,
		key.Size,
		key.ID)

	if err != nil {
		log.Fatal(err)
	}
}
