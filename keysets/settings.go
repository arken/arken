package keysets

import (
	"path/filepath"

	"github.com/arkenproject/arken/ipfs"

	"github.com/arkenproject/arken/config"

	"github.com/BurntSushi/toml"
)

func importKeysetSettings(keyset *config.KeySet, rootPath string) (err error) {
	// Import Settings from keyset.config
	_, err = toml.DecodeFile(filepath.Join(rootPath, "keyset.config"), &keyset)
	if err != nil {
		return err
	}

	err = ipfs.AddBootstrapPeer(keyset.Gateway)
	if err != nil {
		return err
	}

	return nil
}
