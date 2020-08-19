package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"
)

// TransactionCommit adds a transaction to the database table.
func TransactionCommit(tx *sql.Tx, action string, file FileKey) (err error) {
	time.Sleep(1 * time.Second)
	stmt, err := tx.Prepare(
		`INSERT INTO transactions(
			time,
			fileID,
			action,
			size
		) VALUES(datetime('now'),?,?,?);`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(
		file.ID,
		action,
		file.Size)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// TransactionSum returns the sum of all of the transactions for the last month.
func TransactionSum(db *sql.DB, action string) (sum int, err error) {
	// Ping to check that database connection still exists.
	err = db.Ping()
	if err != nil {
		return -1, err
	}
	// Get total pool size from sum of nodes reported values.
	row, err := db.Query("SELECT SUM(size) FROM transactions WHERE action = ? AND time > ?",
		action,
		(time.Now()).AddDate(0, -1, 0))
	if err != nil {
		return -1, err
	}
	defer row.Close()
	if !row.Next() {
		return -1, errors.New("sum not found")
	}
	err = row.Scan(&sum)
	if err != nil && !strings.HasSuffix(err.Error(), "converting NULL to int is unsupported") {
		return -1, err
	}
	return sum, nil
}
