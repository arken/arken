package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/DataDrake/cli-ng/v2/cmd"
	"github.com/arken/arken/config"
	"github.com/inconshreveable/go-update"
	"github.com/tcnksm/go-latest"
)

func init() {
	cmd.Register(&Update)
}

// Update checks for a new version of the Arken program and updates itself
// if a newer version is found and the user agrees to update.
var Update = cmd.Sub{
	Name:  "update",
	Alias: "upd",
	Short: "Update Arken to the latest version available.",
	Args:  &UpdateArgs{},
	Flags: &UpdateFlags{},
	Run:   UpdateRun,
}

// UpdateArgs handles the specific arguments for the update command.
type UpdateArgs struct {
}

// UpdateFlags handles the specific flags for the update command.
type UpdateFlags struct {
	Yes bool `short:"y" long:"yes" desc:"If a newer version is found update without prompting the user."`
}

// UpdateRun handles the checking and self updating of the Arken program.
func UpdateRun(r *cmd.Root, c *cmd.Sub) {
	// Parse Root Flags
	rFlags := r.Flags.(*GlobalFlags)

	fmt.Printf("Current Version: %s\n", config.Version)
	if config.Version == "develop" {
		fmt.Println("Cannot update a development version of Arken.")
		os.Exit(1)
	}

	flags := c.Flags.(*UpdateFlags)
	latestVersion := &latest.GithubTag{
		Owner:      "arken",
		Repository: "arken",
	}

	res, err := latest.Check(latestVersion, config.Version)
	checkError(rFlags, err)

	fmt.Printf("Latest Version: %s\n", res.Current)

	if res.Outdated {
		if !flags.Yes {
			fmt.Println("Would you like to update Arken to the newest version? ([y]/n)")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.ToLower(strings.TrimSpace(input))
			if input == "n" {
				return
			}
		}
		url := "https://github.com/arken/arken/releases/download/v" + res.Current + "/arken-v" + res.Current + "-" + runtime.GOOS + "-" + runtime.GOARCH

		// Display Spinner on Update.
		fmt.Println("Updating Arken...")

		resp, err := http.Get(url)
		checkError(rFlags, err)

		defer resp.Body.Close()
		err = update.Apply(resp.Body, update.Options{})
		checkError(rFlags, err)

		fmt.Print("\rUpdating Arken...Done!\n")
	} else {
		fmt.Println("Already Up-To-Date!")
	}
}
