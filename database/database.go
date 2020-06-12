package database

import (
	"database/sql"

	"github.com/archivalists/arken/config"

	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver for database interaction.
)

// Open opens the Arken SQLite database or creates one if one is not found.
func Open() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", config.Global.Database.Path)
	if err != nil {
		return db, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS keys(
			id TEXT NOT NULL,
			name TEXT NOT NULL,
			size INT(11) NOT NULL,
			status TEXT,

			PRIMARY KEY(id)
		);
			`)
	if err != nil {
		return db, err
	}
	return db, nil
}
