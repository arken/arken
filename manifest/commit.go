package manifest

import (
	"io"
	"os"
	"path/filepath"
)

func (m *Manifest) getCommit() (string, error) {
	// Check if the commit file exists, open it.
	f, err := os.Open(filepath.Join(m.path, "COMMIT"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the commit in
	bytes, err := io.ReadAll(f)
	return string(bytes), err
}

func (m *Manifest) setCommit(commit string) error {
	// Write commit out to file.
	f, err := os.Create(filepath.Join(m.path, "COMMIT"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(commit)
	return err
}
