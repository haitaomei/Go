package main

import (
	"fmt"
	"time"

	"github.com/haitaomei/GoUtil/fileutils"
)

var path = "1.zip"
var newpath = "2.zip"

func main() {
	hashExample(path)
	splitExample()
	hashExample(newpath)
}

func splitExample() {

	chuncks, num, _ := fileutils.Split(path)
	for i := 0; i < len(chuncks); i++ {
		sc, _ := fileutils.DataHash(*chuncks[i])
		fmt.Println("Hash code partition ", i, "= ", sc)
	}
	fileutils.StorePartitions(chuncks, num, path)
	fileutils.Merge(path, num, newpath, true)
}

func hashExample(file string) {
	start := time.Now()
	s, _ := fileutils.FileHash(file)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(file, "'s Hash code", s, "\t using ", elapsed)
}
