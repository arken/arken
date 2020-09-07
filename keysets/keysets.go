package keysets

import (
	"os"
	"path/filepath"

	"github.com/arkenproject/arken/config"

	"github.com/go-git/go-git/v5"
)

// LoadSets takes a list of keysets and indexes the keys within them.
// cloning the repositories if not found locally and pulling updates
// to local repositories.
func LoadSets(keysets []config.KeySet) (err error) {
	for repo := range keysets {
		location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keysets[repo].URL))
		r, err := git.PlainOpen(location)
		if err != nil && err.Error() == "repository does not exist" {
			r, err = git.PlainClone(location, false, &git.CloneOptions{
				URL:               keysets[repo].URL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if err != nil {
				return err
			}
		} else {
			if err != nil {
				return err
			}

			w, err := r.Worktree()
			if err != nil {
				return err
			}

			err = w.Pull(&git.PullOptions{RemoteName: "origin"})
			if err != nil && err.Error() == "already up-to-date" {
				err = importKeysetSettings(&keysets[repo], location)
				if err != nil {
					return err
				}
				continue
			} else if err != nil && err.Error() == "non-fast-forward update" {
				os.RemoveAll(location)
				return LoadSets(keysets)
			} else if err != nil {
				return err
			}
		}
		err = importKeysetSettings(&keysets[repo], location)
		if err != nil {
			return err
		}
	}
	return nil
}
