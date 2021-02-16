package database

import (
	"io"
	"os"

	"github.com/arken/arken/config"
)

// Copy copies the database to a new location.
func Copy(original string, dest string) (err error) {
	// Copy database for long read.
	dbFile, err := os.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}

	copyFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(copyFile, dbFile)
	if err != nil {
		return err
	}
	copyFile.Close()
	dbFile.Close()

	return nil
}
