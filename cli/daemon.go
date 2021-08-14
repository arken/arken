package cli

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/arken/arken/config"
	"github.com/arken/arken/database"
	"github.com/arken/arken/manifest"
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

	// Initialize the node's configuration
	err = config.Init(path)
	checkError(rFlags, err)

	db, err := database.Init(config.Global.Database.Path)
	checkError(rFlags, err)

	// Initialize the node's manifest
	nodeManifest, err := manifest.Init(
		config.Global.Manifest.Path,
		config.Global.Manifest.URL,
	)
	checkError(rFlags, err)

	results, err := nodeManifest.Index(db, false)
	checkError(rFlags, err)

	for result := range results {
		switch result.Status {
		case "add":
			fmt.Println("added:", result)
		case "remove":
			fmt.Println("removed:", result)
		}
	}
}
