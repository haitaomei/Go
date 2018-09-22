package fileutils

import (
	"fmt"
	"io"
	"os"
)

//CreateFile create a new file
func CreateFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}
}

//WriteFile writes a file
func WriteFile(path string, data *[]byte) {
	// open file with RW permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	_, err = file.Write(*data)
	if isError(err) {
		return
	}

	// save changes
	err = file.Sync()
	if isError(err) {
		return
	}
}

//ReadFile read file into a buffer
func ReadFile(path string) ([]byte, error) {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return nil, err
	}
	defer file.Close()

	// read file, line by line
	var data = make([]byte, 1024)
	for {
		_, err = file.Read(data)

		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}
	return data, nil
}

//DeleteFile deletes a file
func DeleteFile(path string) {
	var err = os.Remove(path)
	if isError(err) {
		return
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
