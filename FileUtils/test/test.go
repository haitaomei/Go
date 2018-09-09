package main

import (
	"fmt"
	"time"

	"github.com/haitaomei/GoUtil/fileutils"
)

var path = "/Users/HaiTao/test/1.zip"
var newpath = "/Users/HaiTao/test/2.zip"

func main() {
	hashExample(path)
	splitExample()
	hashExample(newpath)
}

func splitExample() {

	chuncks, _ := fileutils.Split(path)
	for i := 0; i < len(chuncks); i++ {
		sc, _ := fileutils.DataHash(*chuncks[i])
		fmt.Println("Hash code partition ", i, "= ", sc)
	}
	fileutils.StorePartitions(chuncks, path)
	fileutils.Merge(path, len(chuncks), newpath, true)
}

func hashExample(file string) {
	start := time.Now()
	s, _ := fileutils.FileHash(file)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(file, "'s Hash code", s, "\t using ", elapsed)
}
