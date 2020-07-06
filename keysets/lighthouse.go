package keysets

import (
	"errors"
	"path/filepath"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

func configLighthouse(hash string, url string) (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	// Parse URL for Keyset Name
	ksName := filepath.Base(url)

	// Add lighthouse file to database to keep track of last time seen for garbage collector.
	database.Add(db, database.FileKey{ID: hash, Name: "lighthouse", KeySet: ksName})
	if err != nil {
		return errors.New("couldn't add light house key")
	}

	return nil
}
