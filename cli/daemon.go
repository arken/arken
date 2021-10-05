package cli

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/engine"
	"github.com/arken/arken/ipfs"
	"github.com/arken/arken/manifest"
	"github.com/go-co-op/gocron"
)

func init() {
	cmd.Register(&Daemon)
}

var Daemon = cmd.Sub{
	Name:  "daemon",
	Alias: "d",
	Short: "Startup and run the Arken daemon.",
	Run:   RunDaemon,
}

func RunDaemon(r *cmd.Root, s *cmd.Sub) {
	// Parse root flags
	rFlags := r.Flags.(*GlobalFlags)

	// Determine the current user.
	user, err := user.Current()
	checkError(rFlags, err)

	// Determine default program configuration path.
	path := filepath.Join(user.HomeDir, ".arken", "config.toml")

	// Check for custom program configuration path.
	if rFlags.Config != "" {
		path = rFlags.Config
	}

	fmt.Println("Setting up daemon...")

	// Initialize the node's configuration
	err = config.Init(path)
	checkError(rFlags, err)

	db, err := database.Init(config.Global.Database.Path)
	checkError(rFlags, err)

	// Initialize the node's manifest
	manifest, err := manifest.Init(
		config.Global.Manifest.Path,
		config.Global.Manifest.URL,
	)
	checkError(rFlags, err)

	// Initialize embedded IPFS Node
	ipfs, err := ipfs.CreateNode(
		config.Global.Storage.Path,
		ipfs.NodeConfArgs{
			SwarmKey:       manifest.ClusterKey,
			BootstrapPeers: manifest.BootstrapPeers,
			StorageMax:     config.Global.Storage.Limit,
		})
	checkError(rFlags, err)

	// Initialize Arken Engine
	engine := engine.Node{
		Cfg:      &config.Global,
		DB:       db,
		Node:     ipfs,
		Manifest: manifest,
		Verbose:  rFlags.Verbose,
	}
	checkError(rFlags, err)

	// Create Task Scheduler
	tasks := gocron.NewScheduler(time.UTC)

	// Configure Arken Tasks
	// Check for and sync updates to the manifest every hour.
	_, err = tasks.Every(1).Hours().SingletonMode().Do(engine.SyncManifest)
	checkError(rFlags, err)

	// Check the number of times all files in the archive are
	// backed up to determine if any need to be replicated locally
	// to preserve data within the archive.
	rebalance, err := tasks.Every(1).Days().SingletonMode().WaitForSchedule().Do(engine.Rebalance)
	checkError(rFlags, err)
	rebalance.SingletonMode()

	// Verify database consistency against manifest
	verifyDB, err := tasks.Every(5).Days().SingletonMode().WaitForSchedule().Do(engine.VerifyDB)
	checkError(rFlags, err)
	verifyDB.SingletonMode()

	// Very datastore consistency against database
	_, err = tasks.Every(5).Days().SingletonMode().WaitForSchedule().Do(engine.VerifyDatastore)
	checkError(rFlags, err)

	// If stats are enabled send stats to the manifest stats peer.
	if strings.ToLower(config.Global.Stats.Enabled) == "true" {
		fmt.Printf("Stats reporting: enabled\n")
		_, err = tasks.Every(1).Hours().SingletonMode().Do(engine.ReportStats)
		checkError(rFlags, err)
	}

	// Check if verbose mode is enabled.
	if rFlags.Verbose {
		fmt.Printf("Verbose Mode: Enabled\n")
	}

	// Start Task Scheduler
	fmt.Printf("Daemon setup and running\n\n")
	tasks.StartBlocking()

}
