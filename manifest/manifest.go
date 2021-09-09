package manifest

import (
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
)

type Manifest struct {
	Name           string   `toml:"name,omitempty"`
	BootstrapPeers []string `toml:"bootstrap_peers,omitempty"`
	ClusterKey     string   `toml:"cluster_key,omitempty"`
	Replications   int64    `toml:"replications,omitempty"`
	StatsNode      string   `toml:"stats_node,omitempty"`
	url            string   `toml:"url"`
	path           string   `toml:"path"`
	r              *git.Repository
	lock           *sync.Mutex
}

// Init Clones/Pulls a Manifest Repository and Parses the Config
func Init(path, url string) (*Manifest, error) {
	var err error

	// Initialize Manifest Struct
	result := Manifest{
		path: path,
		url:  url,
		lock: &sync.Mutex{},
	}

	// Check if Git Repository Exists
	result.r, err = git.PlainOpen(path)
	if err != nil && err.Error() == "repository does not exist" {
		result.r, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
	}
	if err != nil {
		return nil, err
	}

	// Pull in new changes from the manifest
	err = result.Pull()
	if err != nil {
		return nil, err
	}

	// Decode local manifest config from repository
	err = result.Decode()
	return &result, err
}

// Decode Manifest Configuration from Manifest Repository
func (m *Manifest) Decode() error {
	_, err := toml.DecodeFile(filepath.Join(m.path, "config.toml"), m)
	return err
}
