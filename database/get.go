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
		&result.Status,
		&result.KeySet)
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}

// GetAll opens a channel and reads each entry matching the status into the channel.
func GetAll(db *sql.DB, status string, keySet string, output chan FileKey) {
	err := db.Ping()
	if err != nil {
		close(output)
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM keys WHERE status = ? AND keyset = ?", status, keySet)
	if err != nil {
		close(output)
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var key FileKey

		err = rows.Scan(&key.ID, &key.Name, &key.Size, &key.Status, &key.KeySet)
		if err != nil {
			close(output)
			log.Fatal(err)
		}
		output <- key
	}

	close(output)
}
