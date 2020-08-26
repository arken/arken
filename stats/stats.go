package stats

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/ipfs"
)

// Report uses the uploads stats data to Arken data from the node about
// storage space used, donated, the IPFS identity. Stats can also upload
// optional data like an email to alert users if their nodes
// go offline.
func Report(keysets []config.KeySet) (err error) {
	// Generate Stats Data
	// Parse allotted total amount of storage into bytes.
	total := config.ParseWellFormedPoolSize(config.Global.General.PoolSize)

	// Get the current size of the repo from IPFS in unsigned bytes
	usage, err := ipfs.GetRepoSize()
	if err != nil {
		return err
	}

	// Construct Node Data struct from imported values.
	data := NodeData{
		ID:         ipfs.GetID(),
		Username:   config.Global.Stats.Username,
		Email:      config.Global.Stats.Email,
		TotalSpace: float64(total) / float64(1000000000),
		UsedSpace:  float64(usage) / float64(1000000000),
	}

	for keyset := range keysets {
		for keysets[keyset].LightHouseFileID == "" {
			time.Sleep(30 * time.Second)
		}
		if keysets[keyset].StatsURL != "" {
			fmt.Printf("[Sending Stats Info for: %s]\n", filepath.Base(keysets[keyset].URL))
			err := CheckIn(keysets[keyset].StatsURL, data)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
