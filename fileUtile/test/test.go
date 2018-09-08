package main

import (
	"fmt"
	"time"

	"github.com/haitaomei/GoUtil/fileutile"
)

var path = "1.zip"
var newpath = "2.zip"

func main() {
	hashExample(path)
	splitExample()
	hashExample(newpath)
}

func splitExample() {

	chuncks, num, _ := fileutile.Split(path)
	for i := 0; i < len(chuncks); i++ {
		sc, _ := fileutile.DataHash(*chuncks[i])
		fmt.Println("Hash code partition ", i, "= ", sc)
	}
	fileutile.StorePartitions(chuncks, num, path)
	fileutile.Merge(path, num, newpath, true)
}

func hashExample(file string) {
	start := time.Now()
	s, _ := fileutile.FileHash(file)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(file, "'s Hash code", s, "\t using ", elapsed)
}
