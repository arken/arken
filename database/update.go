package database

import "time"

// Update attempts to modify an existing entry in the database.
func (db *DB) Update(entry File) (old File, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return old, err
	}

	// Check for entry in DB
	old, err = db.get(entry.ID)
	if err != nil {
		return old, err
	}

	// Update the entry if it exists.
	err = db.update(entry)
	return old, err
}

// update changes a file's status in the database.
func (db *DB) update(entry File) (err error) {
	stmt, err := db.conn.Prepare(
		`UPDATE files SET
			name = ?,
			size = ?,
			status = ?,
			replications = ?,
			modified = ?
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		entry.Name,
		entry.Size,
		entry.Status,
		entry.Replications,
		time.Now().UTC(),
		entry.ID,
	)
	return err
}
