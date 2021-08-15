package database

import (
	"database/sql"
	"time"
)

// Get searches for and returns a the corresponding entry from the
// database if the entry exists.
func (db *DB) Get(id string) (result File, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Ping the DB and open a connection if necessary
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	// Get and return entry from DB if it exists
	return db.get(id)
}

// get returns the matching entry from the db if it exists.
func (db *DB) get(id string) (result File, err error) {
	row, err := db.conn.Query("SELECT * FROM files WHERE id = ?", id)
	if err != nil {
		return result, err
	}
	defer row.Close()
	if !row.Next() {
		return result, sql.ErrNoRows
	}
	err = row.Scan(
		&result.ID,
		&result.Name,
		&result.Size,
		&result.Status,
		&result.Modified,
		&result.Replications,
	)

	return result, err
}

func (db *DB) GetAll(status string, limit, page int) (result []File, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Create files slice with limit as size.
	result = []File{}

	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	rows, err := db.conn.Query(
		"SELECT * FROM files WHERE status = ? LIMIT ? OFFSET ?;",
		status,
		limit,
		limit*page,
	)
	if err != nil {
		return result, err
	}

	// Iterate through rows found and insert them into the list.
	for rows.Next() {
		var f File

		err = rows.Scan(
			&f.ID,
			&f.Name,
			&f.Size,
			&f.Status,
			&f.Modified,
			&f.Replications)

		if err != nil {
			rows.Close()
			return nil, err
		}

		result = append(result, f)
	}

	// Check for errors and return
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	if len(result) <= 0 {
		return nil, sql.ErrNoRows
	}

	return result, err
}

func (db *DB) GetAllOlderThan(age time.Time, limit, page int) (result []File, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()

	// Create files slice with limit as size.
	result = make([]File, limit)

	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}

	rows, err := db.conn.Query(
		"SELECT * FROM files WHERE time < ? LIMIT ? OFFSET ?;",
		age,
		limit,
		limit*page,
	)
	if err != nil {
		return result, err
	}

	// Iterate through rows found and insert them into the list.
	for rows.Next() {
		var f File

		err = rows.Scan(
			&f.ID,
			&f.Name,
			&f.Size,
			&f.Status,
			&f.Modified,
			&f.Replications)

		if err != nil {
			rows.Close()
			return nil, err
		}

		result = append(result, f)
	}

	// Check for errors and return
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	if len(result) <= 0 {
		return nil, sql.ErrNoRows
	}

	return result, err
}
