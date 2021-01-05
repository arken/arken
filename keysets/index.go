package keysets

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ipfs/go-cid"
)

// Index extracts the file identifiers from the keyset provided.
func Index(path string, new chan database.FileKey, output chan database.FileKey, settings chan string) (err error) {
	copyName := filepath.Join(filepath.Dir(config.Global.Database.Path), "index.db")
	err = database.Copy(config.Global.Database.Path, copyName)
	defer os.Remove(copyName)

	db, err := database.Open(copyName)
	if err != nil {
		return err
	}
	defer db.Close()

	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	hash, err := database.GetCommit(db)
	if err != nil && err.Error() == "entry not found" {
		err = IndexFull(db, path, new, output)
		if err != nil {
			return err
		}
	} else {
		if ref.Hash().String() != hash {
			hashCommit := plumbing.NewHash(hash)
			err = indexPatch(db, path, hashCommit, new, output)
			if err != nil {
				if err.Error() == "object not found" {
					err = IndexFull(db, path, new, output)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	settings <- ref.Hash().String()

	return nil
}

// IndexFull walks through the repository structure and extracts file identifiers from found
// .ks files.
func IndexFull(db *sql.DB, rootPath string, new chan database.FileKey, output chan database.FileKey) (err error) {
	fileTemplate := database.FileKey{
		Size:   -1,
		Status: "added",
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
					if new != nil {
						new <- fileTemplate
					}
				} else if err == nil {
					output <- entry
				}

			}
			if err := scanner.Err(); err != nil {
				return err
			}
			// Close the file after fully parsed.
			err = file.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func indexPatch(db *sql.DB, path string, commitHash plumbing.Hash, newFiles chan<- database.FileKey, output chan<- database.FileKey) (err error) {
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

	buf := new(bytes.Buffer)
	err = diff.Encode(buf)
	if err != nil {
		return err
	}

	entries := make(map[string]database.FileKey)

	lines := strings.Split(buf.String(), "\n")
	for i := range lines {
		data := strings.Fields(lines[i])
		if len(data) <= 1 {
			continue
		}

		id := strings.TrimPrefix(strings.TrimPrefix(data[0], "+"), "-")
		_, err := cid.Decode(id)
		if err != nil {
			continue
		}

		fileTemplate := database.FileKey{
			Size:   -1,
			Status: "added",
			KeySet: filepath.Base(path)}

		if strings.HasPrefix(lines[i], "+") {

			// Set custom file values.
			fileTemplate.ID = id
			fileTemplate.Name = data[1]

			entry, ok := entries[fileTemplate.ID]
			if !ok {
				entries[fileTemplate.ID] = fileTemplate
			} else if entry.Status == "removed" {
				delete(entries, entry.ID)
			}

		}
		if strings.HasPrefix(lines[i], "-") {

			// Set custom file values.
			fileTemplate.ID = id
			fileTemplate.Name = data[1]
			fileTemplate.Status = "removed"

			entry, ok := entries[fileTemplate.ID]
			if !ok {
				entries[fileTemplate.ID] = fileTemplate
			} else if entry.Status == "added" {
				delete(entries, entry.ID)
			}
		}
	}
	for _, entry := range entries {
		if entry.Status == "added" {
			fmt.Printf("Added: %s  %s\n", entry.ID, entry.Name)
			output <- entry
			_, err := database.Get(db, entry.ID)
			if err != nil && err.Error() == "entry not found" {
				newFiles <- entry
			}
		}
		if entry.Status == "removed" {
			fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
			output <- entry
		}
	}

	return nil
}
