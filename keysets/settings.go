package keysets

import (
	"fmt"
	"path/filepath"

	"github.com/archivalists/arken/config"

	"github.com/BurntSushi/toml"
)

func importKeysetSettings(keyset *config.KeySet, rootPath string) (err error) {
	// Import Settings from keyset.config
	_, err = toml.DecodeFile(filepath.Join(rootPath, "keyset.config"), &keyset)
	if err != nil {
		return err
	}
	fmt.Printf("Found Light House File: %s\n", keyset.LightHouseFileID)
	return nil
}
