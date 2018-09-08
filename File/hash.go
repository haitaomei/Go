package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

func hash(path string) (string, error) {
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

func main() {
	start := time.Now()

	s, _ := hash("hash.go")

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Hash code", s, "\t using ", elapsed)
}
