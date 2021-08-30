package manifest

import (
	"os"

	"github.com/go-git/go-git/v5"
)

// Pull Performs a git pull on the repository
func (m *Manifest) Pull() error {
	// Checkout the repository worktree
	w, err := m.r.Worktree()
	if err != nil {
		return err
	}

	// Save commit before pull
	commit, err := m.getCommit()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check for updates to the Manifest Repository
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	// Write out commit if not nil
	if commit != "" {
		err = m.setCommit(commit)
		if err != nil {
			return err
		}
	}
	return nil
}
