package database

import "time"

// Update attempts to modify an existing entry in the database.
func (db *DB) Update(entry File) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return err
	}

	// Update the entry if it exists.
	return db.update(entry)
}

// update changes a file's status in the database.
func (db *DB) update(entry File) (err error) {
	stmt, err := db.conn.Prepare(
		`UPDATE files SET
			status = ?,
			replications = ?,
			size = ?,
			modified = ?,
			used_space = ?
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.Status,
		entry.Replications,
		entry.Size,
		time.Now().UTC(),
		entry.ID,
	)
	return err
}
