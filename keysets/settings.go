package keysets

import (
	"path/filepath"
	"strings"

	"github.com/arken/arken/config"

	"github.com/BurntSushi/toml"
)

func importKeysetSettings(keyset *config.KeySet, rootPath string) (err error) {
	// Import Settings from keyset.config
	_, err = toml.DecodeFile(filepath.Join(rootPath, "keyset.config"), keyset)
	if err != nil {
		return err
	}
	keyset.Name = strings.TrimSuffix(filepath.Base(rootPath), ".git")
	keyset.Path = rootPath

	return nil
}
