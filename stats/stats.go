package stats

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/alecthomas/units"
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
	total, err := units.ParseStrictBytes(config.Global.General.PoolSize)
	if err != nil {
		return err
	}

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
		TotalSpace: float64(total / 1000000000),
		UsedSpace:  float64(usage / 1000000000),
	}

	for keyset := range keysets {
		if keysets[keyset].StatsURL != "" {
			fmt.Printf("[Sending Stats Info for: %s]\n", filepath.Base(keysets[keyset].URL))
			err := CheckIn(keysets[keyset].StatsURL, data)
			if err != nil {
				log.Println(err)
			}
		}
	}

	fmt.Println(data)
	return nil
}
