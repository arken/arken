package database

// Remove deletes and returns an entry from the database.
func (db *DB) Remove(id string) (result File, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	// Get the current value of the entry in the DB before removing
	result, err = db.get(id)
	if err != nil {
		return result, err
	}

	// Remove the entry from the DB
	err = db.remove(id)
	return result, err
}

// remove deletes an entry to the DB.
func (db *DB) remove(id string) (err error) {
	stmt, err := db.conn.Prepare(
		"DELETE FROM files WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}
