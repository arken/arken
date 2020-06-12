package keysets

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func index(rootPath string) {
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(filepath.Base(path), ".ks") {
			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}
