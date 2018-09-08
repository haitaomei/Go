package fileutils

import (
	"os"
	"path/filepath"
)

func filterFiles(files *[]string, includeDir bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fi, _ := os.Stat(path)
		if !fi.IsDir() {
			*files = append(*files, path)
		} else {
			if includeDir {
				*files = append(*files, path)
			}
		}
		return nil
	}
}

// ListFiles get all the files in a directory
func ListFiles(path string) (*[]string, error) {
	files, err := enum(path, false)
	return files, err
}

// ListFilesDir get all the files and directories in a directory
func ListFilesDir(path string) (*[]string, error) {
	files, err := enum(path, true)
	return files, err
}

func enum(path string, includeDir bool) (*[]string, error) {
	var files []string
	err := filepath.Walk(path, filterFiles(&files, includeDir))
	if err != nil {
		return nil, err
	}
	return &files, nil
}
