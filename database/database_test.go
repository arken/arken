package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import sqlite driver for database interaction.
)

// openMock opens up a mock DB for unit testing purposes.
func openMock() (result *DB, err error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
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

	result = &DB{
		conn: db,
		lock: sync.Mutex{},
	}

	return result, err

}
