package keysets

import (
	"path/filepath"

	"github.com/arkenproject/arken/database"
)

// ConfigLighthouse constructs a lighthouse file key.
func ConfigLighthouse(hash string, url string) (result database.FileKey, err error) {
	// Parse URL for Keyset Name
	ksName := filepath.Base(url)

	return database.FileKey{
		ID:     hash,
		Name:   "lighthouse",
		KeySet: ksName}, nil
}
