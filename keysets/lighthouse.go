package keysets

import (
	"errors"
	"fmt"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

func configLighthouse(hash string) (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	// Add lighthouse file to database to keep track of last time seen for garbage collector.
	database.Add(db, database.FileKey{ID: hash, Name: "lighthouse"})
	if err != nil {
		return errors.New("couldn't add light house key")
	}

	// Pin Lighthouse File to determine the size of the active cluster.
	fmt.Println("Pinning Lighthouse File...")
	err = ipfs.Pin(hash)
	if err != nil {
		return err
	}
	return nil
}
