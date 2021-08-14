package manifest

import (
	"bufio"
	"bytes"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arken/arken/database"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ipfs/go-cid"
)

func (m *Manifest) Index(db *database.DB, force bool) (<-chan database.File, error) {
	// Check to see if an existing commit checkpoint exists
	commit, err := m.getCommit()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Grab the current commit of the repository
	ref, err := m.r.Head()
	if err != nil {
		return nil, err
	}

	// Create results channel
	results := make(chan database.File, 50)

	// Initialize go routine to handle indexing
	go func() {
		// Check if fullIndex or diffIndex should be used.
		switch {
		case commit != "" && ref.Hash().String() != commit && !force:
			m.indexDiff(commit, results)
		case commit == "" || force:
			m.indexFull(db, results)
		}

		// Save current git hash as checkpoint.
		m.setCommit(ref.Hash().String())

		// Close results channel
		close(results)
	}()

	return results, nil
}

func (m *Manifest) indexFull(db *database.DB, results chan<- database.File) {
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

				// Check if entry is already in the database.
				if _, err := db.Get(data[0]); err == nil || err != sql.ErrNoRows {
					continue
				}

				// Only add files not in the DB to the added channel
				results <- database.File{
					ID:     data[0],
					Name:   data[1],
					Status: "add",
					Size:   0,
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
	if err != nil {
		log.Println(err)
	}

	// Check for deleted files by looking for everything not touched by the index
	for i := 0; ; i++ {
		files, err := db.GetAllOlderThan(start, 100, i)
		if err != nil {
			break
		}
		for _, file := range files {
			file.Status = "remove"
			results <- file
		}
	}
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
}

func (m *Manifest) indexDiff(oldCommit string, results chan<- database.File) {
	r, err := git.PlainOpen(m.path)
	if err != nil {
		log.Println(err)
		return
	}
	ref, err := r.Head()
	if err != nil {
		log.Println(err)
		return
	}
	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Println(err)
		return
	}

	parent, err := r.CommitObject(plumbing.NewHash(oldCommit))
	if err != nil {
		log.Println(err)
		return
	}

	diff, err := parent.Patch(commit)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = diff.Encode(buf)
	if err != nil {
		log.Println(err)
		return
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

		fileTemplate := database.File{}

		if strings.HasPrefix(lines[i], "+") {

			// Set custom file values.
			fileTemplate.ID = id
			fileTemplate.Name = data[1]
			fileTemplate.Status = "add"

			entry, ok := entries[fileTemplate.ID]
			if !ok {
				entries[fileTemplate.ID] = fileTemplate
			} else if entry.Status == "remove" {
				delete(entries, entry.ID)
			}

		}
		if strings.HasPrefix(lines[i], "-") {

			// Set custom file values.
			fileTemplate.ID = id
			fileTemplate.Name = data[1]
			fileTemplate.Status = "remove"

			entry, ok := entries[fileTemplate.ID]
			if !ok {
				entries[fileTemplate.ID] = fileTemplate
			} else if entry.Status == "add" {
				delete(entries, entry.ID)
			}
		}
	}
	for _, entry := range entries {
		if entry.Status == "add" {
			results <- entry
		}
		if entry.Status == "remove" {
			results <- entry
		}
	}
}
