package manifest

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arken/arken/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ipfs/go-cid"
)

// IndexOptions lists the input/output options for the index functions.
type IndexOptions struct {
	DB      *database.DB
	Added   chan<- database.File
	Removed chan<- database.File
	Errors  chan<- error
	Verbose bool
}

func (m *Manifest) Index(opts IndexOptions) {
	// Check to see if an existing commit checkpoint exists
	commit, err := m.getCommit()
	if err != nil && !os.IsNotExist(err) {
		opts.Errors <- err
	}

	r, err := git.PlainOpen(m.path)
	if err != nil {
		opts.Errors <- err
	}
	ref, err := r.Head()
	if err != nil {
		opts.Errors <- err
	}

	switch {
	case commit != "" && ref.Hash().String() != commit:
		m.indexDiff(opts)
	case commit == "":
		m.IndexFull(opts)
	}
}

func (m *Manifest) IndexFull(opts IndexOptions) {
	// Construct a file template for new files.
	fileTemplate := database.File{
		Size:   -1,
		Status: "remote",
	}

	// Store current time to find files not touched during index
	start := time.Now().UTC()

	// Walk through entire repository directory structure to look for .ks files.
	err := filepath.Walk(m.path, func(path string, info os.FileInfo, err error) error {
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
				if opts.Verbose {
					fmt.Printf("Parsed: %s\n", data)
				}

				// Set custom file values.
				fileTemplate.ID = data[0]
				fileTemplate.Name = data[1]

				// Check if entry is already in the database.
				if _, err := opts.DB.Get(fileTemplate.ID); err != nil &&
					err == sql.ErrNoRows {
					continue
				}

				// Only add files not in the DB to the added channel
				opts.Added <- fileTemplate
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
	if err != nil {
		opts.Errors <- err
	}

	// Check for deleted files by looking for everything not touched by the index
	for i := 0; ; i++ {
		files, err := opts.DB.GetAllOlderThan(start, 100, i)
		if err != nil {
			break
		}
		for _, file := range files {
			opts.Removed <- file
		}
	}
	if err != nil && err != sql.ErrNoRows {
		opts.Errors <- err
	}
	opts.Errors <- nil
}

func (m *Manifest) indexDiff(opts IndexOptions) {
	oldCommit, err := m.getCommit()
	if err != nil {
		opts.Errors <- err
	}

	r, err := git.PlainOpen(m.path)
	if err != nil {
		opts.Errors <- err
	}
	ref, err := r.Head()
	if err != nil {
		opts.Errors <- err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		opts.Errors <- err
	}

	parent, err := r.CommitObject(plumbing.NewHash(oldCommit))
	if err != nil {
		opts.Errors <- err
	}

	diff, err := parent.Patch(commit)
	if err != nil {
		opts.Errors <- err
	}

	buf := new(bytes.Buffer)
	err = diff.Encode(buf)
	if err != nil {
		opts.Errors <- err
	}

	entries := make(map[string]database.File)

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

		fileTemplate := database.File{
			Size:   -1,
			Status: "added",
		}

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
			if opts.Verbose {
				fmt.Printf("Added: %s  %s\n", entry.ID, entry.Name)
			}
			entry.Status = "remote"
			opts.Added <- entry
		}
		if entry.Status == "removed" {
			if opts.Verbose {
				fmt.Printf("Removed: %s  %s\n", entry.ID, entry.Name)
			}
			opts.Removed <- entry
		}
	}

	err = m.setCommit(commit.Hash.String())
	opts.Errors <- err
}
