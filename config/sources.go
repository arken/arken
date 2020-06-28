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
	Keysets  []KeySet
	internal sourcesFileData
)

// KeySet defines the global structure housing the known keysets.
type KeySet struct {
	URL               string
	LightHouseFileID  string
	ReplicationFactor float32
	Gateway           string
}

type sourcesFileData struct {
	Sets []string `yaml:"keysets"`
}

func readSources() {
	// Parse Keysets Yaml file.
	fileData, err := ioutil.ReadFile(Global.Sources.Config)
	if os.IsNotExist(err) {
		genSources(defaultSources())
		readSources()
		return
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(fileData, &internal)
	if err != nil {
		log.Fatal(err)
	}

	for set := range internal.Sets {
		Keysets = append(Keysets, KeySet{URL: internal.Sets[set]})
	}
}

func defaultSources() sourcesFileData {
	result := sourcesFileData{
		Sets: []string{"https://github.com/arkenproject/core-keyset"},
	}

	return result
}

func genSources(keyset sourcesFileData) {
	os.MkdirAll(filepath.Dir(Global.Sources.Config), os.ModePerm)

	out, err := yaml.Marshal(keyset)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filepath.Join(filepath.Dir(Global.Sources.Config), "keysets.yaml"), out, 0664)
	if err != nil {
		log.Fatal(err)
	}
}
