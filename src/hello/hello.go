package main

import (
	"GoPractice/src/stringutil"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(stringutil.Reverse("!oG ,olleH"))
	fmt.Println("My favorite number is", rand.Intn(10), "--")
	fmt.Println(add(3, 4))
	fmt.Println(swap("World", "Hello "))
}

func add(x int, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}
