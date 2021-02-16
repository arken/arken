package tasks

import (
	"log"
	"strings"

	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/engine"
	"github.com/arken/arken/keysets"
)

// Main handles the main parent loop of Arken's Management System.
func Main() {
	// New is reserved for incomming files from the indexer.
	new := make(chan database.FileKey, 10)
	// Remove is for existing remote files from the database.
	remote := make(chan database.FileKey, 10)
	// Contents of output will be added to the database.
	output := make(chan database.FileKey, 20)
	// Contents of settings will be hash strings for checkpointing database.
	settings := make(chan string, 1)

	// Initialize Database Writer Function
	go databaseWriter(output, settings)

	// Load KeySet Configurations on Boot
	err := keysets.LoadSets(config.Keysets)
	if err != nil {
		log.Fatal(err)
	}

	// Check the Size of the Node Storage Utilization
	err = checkNodeSize(output)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Keyset Refresh Task
	go loadSets(config.Keysets, new, output, settings)

	// Initialize Engine Network Limit Test
	go checkNetworkLimit()

	// Initialize Database Reader
	go databaseReader(remote, output)

	// Initialize Database-Keyset Verifacation Test
	go verifyDB(config.Keysets, new, output)

	// Verify Locally Pinned Files and Re-Pin if lost.
	go VerifyLocal()

	// Launch Stats Reporting if enabled in the config.
	if strings.ToLower(config.Global.General.StatsReporting) == "on" {
		go StatsReporting()
	}

	err = engine.Run(new, remote, output)
	if err != nil {
		log.Fatal(err)
	}
}
