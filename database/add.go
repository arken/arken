package database

import (
	"database/sql"
	"time"
)

// Add inserts a File entry into the database if it doesn't exist already.
func (db *DB) Add(input File) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return err
	}

	_, err = db.get(input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.insert(input)
		} else {
			return err
		}
	} else {
		err = db.update(input)
	}
	return err
}

// Insert adds a Node entry to the database.
func (db *DB) insert(entry File) (err error) {
	stmt, err := db.conn.Prepare(
		`INSERT INTO files(
			id,
			name,
			size,
			status,
			modified,
			replications
		) VALUES(?,?,?,?,?,?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.ID,
		entry.Name,
		entry.Size,
		entry.Status,
		time.Now().UTC(),
		entry.Replications,
	)
	return err
}
