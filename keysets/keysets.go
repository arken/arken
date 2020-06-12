package keysets

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/archivalists/arken/config"

	"github.com/go-git/go-git/v5"
)

// Index takes a list of keysets and indexes the keys within them.
// cloning the repositories if not found locally and pulling updates
// to local repositories.
func Index(keysets []string) {
	for repo := range keysets {
		fmt.Println(keysets[repo])
		location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keysets[repo]))
		r, err := git.PlainOpen(location)
		if err != nil && err.Error() == "repository does not exist" {
			r, err = git.PlainClone(location, false, &git.CloneOptions{
				URL:               keysets[repo],
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if err != nil {
				log.Fatal(err)
			}

		} else {
			if err != nil {
				log.Fatal(err)
			}

			w, err := r.Worktree()
			if err != nil {
				log.Fatal(err)
			}

			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err.Error() != "already up-to-date" {
				log.Fatal(err)
			}
		}

	}
}
