package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	// Keysets is the configuration structure for the known keysets.
	Keysets KeySources
)

// KeySources defines the global structure housing the known keysets.
type KeySources struct {
	Sets []string `yaml:"keysets"`
}

func readSources(keysets *KeySources) {
	// Parse Keysets Yaml file.
	fileData, err := ioutil.ReadFile(Global.Sources.Config)
	if os.IsNotExist(err) {
		genSources(defaultSources())
		readSources(keysets)
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(fileData, &Keysets)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultSources() KeySources {
	result := KeySources{
		Sets: []string{"https://github.com/archivalists/core-keyset-testing"},
	}

	return result
}

func genSources(keyset KeySources) {
	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	out, err := yaml.Marshal(keyset)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(filepath.Dir(path), "keysets.yaml"), out, 0664)
	if err != nil {
		log.Fatal(err)
	}
}
