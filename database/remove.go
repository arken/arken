package database

import (
	"database/sql"
	"log"
)

// Delete removes an entry from the database.
func Delete(db *sql.DB, id string) error {
	stmt, err := db.Prepare(
		"DELETE FROM keys WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	return err
}
