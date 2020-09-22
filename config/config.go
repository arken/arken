package config

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config defines the configuration struct for importing settings from TOML.
type Config struct {
	General  general
	Database database
	Sources  sources
	Stats    stats
}

// general defines the substruct about general application settings.
type general struct {
	Version        string
	PoolSize       string
	NetworkLimit   string
	StatsReporting string
}

// database defines database specific config settings.
type database struct {
	Path string
}

// sources defines where to look for the local cloned Keyset repositories
type sources struct {
	Config       string
	Repositories string
	Storage      string
}

// stats defines where to look for the stats configuration.
type stats struct {
	Username string
	Email    string
}

var (
	// Global is the configuration struct for the application.
	Global Config
	// Disk is the configuration interface for the disk utilities.
	Disk DiskInfo
	path string
)

// initialize the app config system. If a config doesn't exist, create one.
// If the config is out of date read the current config and rebuild with new fields.
func init() {
	// Determine the current user to build expected file path.
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	// Create expected config path.
	path = filepath.Join(user.HomeDir, ".arken", "arken.config")
	readConf(&Global)
	// If the configuration version has changed update the config to the new
	// format while keeping the user's preferences.
	if Global.General.Version != defaultConf().General.Version {
		reloadConf()
		readConf(&Global)
	}
	ConsolidateEnvVars(&Global)
	readSources()

	err = createSwarmKey()
	if err != nil {
		log.Fatal(err)
	}
}

// LoadDiskConfig loads the Disk Configuration
func LoadDiskConfig() {
	ParsePoolSize(&Disk)
	Global.General.PoolSize = Disk.GetPrettyPoolSize()
}

// Read the config or create a new one if it doesn't exist.
func readConf(conf *Config) {
	_, err := toml.DecodeFile(path, &conf)
	if os.IsNotExist(err) {
		genConf(defaultConf())
		readConf(conf)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func createSwarmKey() (err error) {
	keyData := []byte(`/key/swarm/psk/1.0.0/
/base16/
793bdb68b7cfd2f49071a299711df51f1c60283a047e4a8756a5c3a3d1ab776f`)

	os.MkdirAll(Global.Sources.Storage, os.ModePerm)
	err = ioutil.WriteFile(filepath.Join(Global.Sources.Storage, "swarm.key"), keyData, 0644)
	return err
}
