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

func index(rootPath string) (err error) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		return err
	}
	defer db.Close()

	// Walk through entire repository directory structure to look for .ks files.
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// On each interation of the "walk" this function will check if a keyset
		// file and parse for file IDs if true.
		if strings.HasSuffix(filepath.Base(path), ".ks") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				data := strings.Split(scanner.Text(), " : ")
				fmt.Println(data)
				err = database.Add(db, database.FileKey{ID: data[1], Name: data[0], Size: -1, Status: "remote", KeySet: filepath.Base(rootPath)})
				if err != nil {
					return err
				}
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			file.Close()
		}
		return nil
	})
	return nil
}
