package config

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/disk"
)

var (
	// Version is the current version of Arken
	Version string = "develop"
	// Global is the global application configuration
	Global config
)

type config struct {
	Database database `toml:"database"`
	Storage  storage  `toml:"storage"`
	Manifest manifest `toml:"manifest"`
	Network  network  `toml:"network"`
	Stats    stats    `toml:"stats"`
}

type database struct {
	Path string `toml:"path"`
}

type storage struct {
	Limit string `toml:"limit"`
	Path  string `toml:"path"`
}

type manifest struct {
	Name           string `toml:"name,omitempty"`
	BootstrapPeers string `toml:"bootstrap_peers,omitempty"`
	ClusterKey     string `toml:"cluster_key,omitempty"`
	Replications   string `toml:"replications,omitempty"`
	StatsNode      string `toml:"stats_node,omitempty"`
	URL            string `toml:"url"`
}

type network struct {
	Limit string `toml:"limit"`
}

type stats struct {
	Enabled string `toml:"enabled"`
	Email   string `toml:"email"`
}

func Init(path string) error {

	// Generate the default config
	Global = config{
		Database: database{
			Path: filepath.Join(filepath.Dir(path), "arken.db"),
		},
		Storage: storage{
			Limit: "50GB",
			Path:  filepath.Join(filepath.Dir(path), "storage"),
		},
		Manifest: manifest{
			URL: "https://github.com/arken/core-manifest",
		},
		Network: network{
			Limit: "500GB",
		},
		Stats: stats{
			Enabled: "true",
			Email:   "",
		},
	}

	// Read in config from file
	err := parseFile(path, &Global)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Read in config from environment
	err = sourceEnv(&Global)
	if err != nil {
		return err
	}

	// As long as the IPFS nodes hasn't been initialized set the storage quota
	// to the max storage available on the drive.
	if _, err := os.Open(filepath.Join(Global.Storage.Path, "config")); os.IsNotExist(err) {

		// Create the storage path if it doesn't exist.
		if _, err := os.Open(Global.Storage.Path); os.IsNotExist(err) {
			err = os.MkdirAll(Global.Storage.Path, os.ModePerm)
			if err != nil {
				return err
			}
		}

		// Grab Disk Usage Stats
		dStat, err := disk.Usage(Global.Storage.Path)
		if err != nil {
			return err
		}

		// Get Disk Size from Config
		cSize, err := humanize.ParseBytes(Global.Storage.Limit)
		if err != nil {
			return err
		}

		// Check if the detected disk size is smaller than the disk limit
		if dStat != nil && cSize > dStat.Free {
			Global.Storage.Limit = humanize.Bytes(dStat.Free)
		}
	}

	// Write config file
	err = writeFile(path, &Global)
	return err
}

func parseFile(path string, in *config) error {
	_, err := toml.DecodeFile(path, in)
	return err
}

func sourceEnv(in *config) error {
	numSubStructs := reflect.ValueOf(in).Elem().NumField()
	// Check for env args matching each of the sub structs.
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(in).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())
		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			fieldName := structType.Field(j).Name
			evName := "ARKEN" + "_" + subStruct + "_" + strings.ToUpper(fieldName)
			evVal, evExists := os.LookupEnv(evName)
			if evExists && evVal != fieldVal {
				iter.FieldByName(fieldName).SetString(evVal)
			}
		}
	}
	return nil
}

func writeFile(path string, in *config) error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(in)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, buf.Bytes(), os.ModePerm)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, buf.Bytes(), os.ModePerm)
	}
	return err
}
