package keysets

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/arkenproject/arken/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// LoadSets takes a list of keysets and indexes the keys within them.
// cloning the repositories if not found locally and pulling updates
// to local repositories.
func LoadSets(keysets []config.KeySet) (err error) {
	newClient := &http.Client{Timeout: 5 * time.Second}

	// Override http(s) default protocol to use our custom client
	client.InstallProtocol("https", githttp.NewClient(newClient))

	for repo := range keysets {
		location := filepath.Join(config.Global.Sources.Repositories, filepath.Base(keysets[repo].URL))
		r, err := git.PlainOpen(location)
		if err != nil && err.Error() == "repository does not exist" {
			r, err = git.PlainClone(location, false, &git.CloneOptions{
				URL:               keysets[repo].URL,
				RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			})

			if err != nil {
				newClient.CloseIdleConnections()
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
				newClient.CloseIdleConnections()
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
