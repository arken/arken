package engine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/ipfs"
)

// Rebalance manages balancing new and at risk files
// between nodes.
func Rebalance() (err error) {
	// Copy database for long read.
	dbFile, err := os.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}

	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "copy.db")
	copyFile, err := os.Create(copyName)
	if err != nil {
		return err
	}
	defer os.Remove(copyName)

	_, err = io.Copy(copyFile, dbFile)
	if err != nil {
		return err
	}
	copyFile.Close()
	dbFile.Close()

	// Open Database and Copy

	write, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	defer write.Close()

	read, err := database.Open(copyName)
	if err != nil {
		return err
	}
	defer read.Close()

	for set := range config.Keysets {
		// Pin Lighthouse File to determine the size of the active cluster.
		fmt.Println("Pinning Lighthouse File...")
		err = ipfs.Pin(config.Keysets[set].LightHouseFileID)
		if err != nil {
			return err
		}

		keySet := filepath.Base(config.Keysets[set].URL)

		threshold, err := CalcThreshold(config.Keysets[set].LightHouseFileID, config.Keysets[set].Replications, 20)
		if err != nil {
			return err
		}

		fmt.Printf("Found Threshold To Be: %d\n", threshold)

		err = ScanReplications(read, write, keySet, threshold)
		if err != nil {
			return err
		}
	}

	return nil
}
