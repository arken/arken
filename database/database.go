package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver for database interaction.
)

// FileKey is a struct format of data within the Keys database.
type FileKey struct {
	ID     string
	Name   string
	Size   int
	Status string
}

// Open opens the Arken SQLite database or creates one if one is not found.
func Open(path string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS keys(
			id TEXT NOT NULL,
			name TEXT NOT NULL,
			size INT(11),
			status TEXT NOT NULL,

			PRIMARY KEY(id)
		);
			`)
	if err != nil {
		return db, err
	}
	return db, nil
}
