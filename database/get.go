package database

import (
	"database/sql"
	"errors"
	"log"
)

// Get searches for and returns a the coorisponding entry from the
// database if the entry exists.
func Get(db *sql.DB, id string) (result FileKey, err error) {
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	row, err := db.Query("SELECT * FROM keys WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		return result, errors.New("entry not found")
	}
	err = row.Scan(
		&result.ID,
		&result.Name,
		&result.Size,
		&result.Status)
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}
