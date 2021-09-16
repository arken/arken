package database

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import sqlite driver for database interaction.
)

// DB is a wrapper for the database.
type DB struct {
	conn *sql.DB
	lock sync.Mutex
}

// File is a struct of info for a file within Arken.
type File struct {
	ID           string
	Name         string
	Size         int64
	Status       string
	Modified     time.Time
	Replications int
}

// Init opens and connects to the database.
func Init(path string) (result *DB, err error) {
	db, err := sql.Open("sqlite3", path+"?cache=shared")
	if err != nil {
		return nil, err
	}

	// Create the nodes table if is doesn't already exist.
	// This will also create the database if it doesn't exist.
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS files(
			id TEXT NOT NULL,
			name TEXT,
			size INT(11),
			status TEXT NOT NULL,
			modified DATETIME,
			replications INT,
			PRIMARY KEY(id)
		);`,
	)
	if err != nil {
		return nil, err
	}

	result = &DB{
		conn: db,
		lock: sync.Mutex{},
	}
	return result, nil
}
