package keysets

import (
	"path/filepath"

	"github.com/arken/arken/database"
	"github.com/arken/arken/ipfs"
)

// ConfigLighthouse constructs a lighthouse file key.
func ConfigLighthouse(hash string, url string) (result database.FileKey, err error) {
	// Parse URL for Keyset Name
	ksName := filepath.Base(url)
	// Pin Lighthouse file.
	err = ipfs.Pin(hash)
	if err != nil {
		return result, err
	}

	return database.FileKey{
		ID:     hash,
		Status: "local",
		Name:   "lighthouse",
		KeySet: ksName}, nil
}
