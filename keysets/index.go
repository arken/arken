package keysets

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arkenproject/arken/config"
	"github.com/arkenproject/arken/database"
)

// index walks through the repository structure and extracts file identifiers from found
// .ks files.
func index(rootPath string) (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	// Wait to close the database until all files have been indexed.
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

				// Add parsed file to database.
				err = database.Add(db, fileTemplate)
				if err != nil {
					return err
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
