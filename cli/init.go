package cli

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/arken/arken/config"
)

func init() {
	cmd.Register(&Init)
}

var Init = cmd.Sub{
	Name:  "init",
	Alias: "i",
	Short: "Initialize Arken's configuration without running the daemon.",
	Run:   RunInit,
}

func RunInit(r *cmd.Root, s *cmd.Sub) {
	// Parse Root Flags
	rFlags := r.Flags.(*GlobalFlags)

	// Determine Program Configuration Path
	var path string
	if rFlags.Config != "" {
		path = rFlags.Config
	} else {
		user, err := user.Current()
		checkError(rFlags, err)
		path = filepath.Join(user.HomeDir, ".arken", "config.toml")
	}

	fmt.Printf("Initializing configuration at:\n  %s\n\n", path)

	// Initialize Arken Config
	err := config.Init(path)
	checkError(rFlags, err)

	fmt.Printf("Done!\n")

}
