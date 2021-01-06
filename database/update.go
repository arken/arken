package database

import (
	"database/sql"
	"time"
)

// Update changes a file's status in the database.
func Update(db *sql.DB, key FileKey) error {
	stmt, err := db.Prepare(
		`UPDATE keys SET
			Status = ?,
			Replications = ?,
			Size = ?,
			Modified = ?
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		key.Status,
		key.Replications,
		key.Size,
		time.Now().UTC(),
		key.ID)

	return err
}
