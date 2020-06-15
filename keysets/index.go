package keysets

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/archivalists/arken/config"
	"github.com/archivalists/arken/database"
)

func index(rootPath string) {
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(filepath.Base(path), ".ks") {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				data := strings.Split(scanner.Text(), " : ")
				fmt.Println(data)
				err = database.Add(db, database.FileKey{ID: data[1], Name: data[0], Size: -1, Status: "remote"})
				if err != nil {
					log.Fatal(err)
				}
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}
