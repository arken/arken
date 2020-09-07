package keysets

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Index extracts the file identifiers from the keyset provided.
func Index(path string, new chan database.FileKey, output chan database.FileKey) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	if config.Global.General.IndexHash == "" {
		err = indexFull(path, new, output)
		if err != nil {
			return err
		}
	} else {
		if ref.Hash().String() != config.Global.General.IndexHash {
			hash := plumbing.NewHash(config.Global.General.IndexHash)
			err = indexPatch(path, hash, new, output)
			if err != nil {
				return err
			}
		}
	}
	config.Global.General.IndexHash = ref.Hash().String()
	config.GenConf(config.Global)

	return nil
}

// IndexFull walks through the repository structure and extracts file identifiers from found
// .ks files.
func indexFull(rootPath string, new chan database.FileKey, output chan database.FileKey) (err error) {
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "index.db")
	err = database.Copy(config.Global.Database.Path, copyName)
	defer os.Remove(copyName)

	db, err := database.Open(copyName)
	if err != nil {
		return err
	}
	defer db.Close()

	fileTemplate := database.FileKey{
		Size:   -1,
		Status: "remote",
		KeySet: filepath.Base(rootPath)}

	// Walk through entire repository directory structure to look for .ks files.
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// On each interation of the "walk" this function will check if a keyset
		// file and parse for file IDs if true.
		if strings.HasSuffix(filepath.Base(path), ".ks") {
			file, err := os.Open(path)
			// This shouldn't be an error unless part of the keyset was
			// corrupted in transit.
			if err != nil {
				return err
			}
			// Open the files for reading.
			scanner := bufio.NewScanner(file)
			// Scan one line at a time.
			for scanner.Scan() {
				// Split data on white space.
				data := strings.Fields(scanner.Text())
				fmt.Printf("Parsed: %s\n", data)

				// Set custom file values.
				fileTemplate.ID = data[0]
				fileTemplate.Name = data[1]

				entry, err := database.Get(db, data[0])
				if err != nil && err.Error() == "entry not found" {
					// Send parsed file to engine.
					output <- fileTemplate
					new <- fileTemplate
				} else if err == nil {
					output <- entry
				}

			}
			if err := scanner.Err(); err != nil {
				return err
			}
			// Close the file after fully parsed.
			file.Close()
		}
		return nil
	})
	return err
}

func indexPatch(path string, commitHash plumbing.Hash, new chan<- database.FileKey, output chan<- database.FileKey) (err error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	parent, err := r.CommitObject(commitHash)
	if err != nil {
		return err
	}

	diff, err := parent.Patch(commit)
	if err != nil {
		return err
	}

	fileTemplate := database.FileKey{
		Size:   -1,
		Status: "remote",
		KeySet: filepath.Base(path)}

	lines := strings.Split(diff.String(), "\n")
	for i := range lines {
		if strings.HasPrefix(lines[i], "+Qm") {
			data := strings.Fields(lines[i])
			fmt.Printf("Added: %s\n", data)

			// Set custom file values.
			fileTemplate.ID = strings.TrimPrefix(data[0], "+")
			fileTemplate.Name = data[1]
			output <- fileTemplate
			new <- fileTemplate
		}
		if strings.HasPrefix(lines[i], "-Qm") {
			data := strings.Fields(lines[i])
			fmt.Printf("Removed: %s\n", data)

			// Set custom file values.
			fileTemplate.ID = strings.TrimPrefix(data[0], "-")
			fileTemplate.Name = data[1]
			fileTemplate.Status = "removed"
			output <- fileTemplate
		}

	}
	return nil
}
