package fileutils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// FileHash calculate SHA256 code of a file
func FileHash(path string) (string, error) {
	h := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s, nil
}

// DataHash calculate SHA256 code of a data
func DataHash(data *[]byte) (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256(*data)), nil
}
