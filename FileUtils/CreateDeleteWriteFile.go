package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
)

//CreateFile create a new file, if directory doesn't exist, create one
func CreateFile(path string) {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, os.ModePerm)

	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if checkErr(err) {
			return
		}
		defer file.Close()
	}
}

//WriteFile writes a file
func WriteFile(path string, data *[]byte) {
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if checkErr(err) {
		return
	}
	defer file.Close()

	_, err = file.Write(*data)
	if checkErr(err) {
		return
	}
	err = file.Sync()
	if checkErr(err) {
		return
	}
}

//DeleteFile deletes a file
func DeleteFile(path string) {
	var err = os.Remove(path)
	if checkErr(err) {
		return
	}
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}
