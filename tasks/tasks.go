package tasks

import (
	"log"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/arkenproject/arken/engine"
	"github.com/arkenproject/arken/keysets"
)

// Main handles the main parent loop of Arken's Management System.
func Main() {
	// New is reserved for incomming files from the indexer.
	new := make(chan database.FileKey, 10)
	// Remove is for existing remote files from the database.
	remote := make(chan database.FileKey, 10)
	// Contents of output will be added to the database.
	output := make(chan database.FileKey, 10)

	// Open connection to Database
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Database Writer Function
	go databaseWriter(db, output)

	// Load KeySet Configurations on Boot
	err = keysets.LoadSets(config.Keysets)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Keyset Refresh Task
	go loadSets(config.Keysets, new, output)

	// Initialize Engine Network Limit Test
	go checkNetworkLimit()

	// Initialize Database Reader
	go databaseReader(remote)

	err = engine.Run(new, remote, output)
	if err != nil {
		log.Fatal(err)
	}
}
