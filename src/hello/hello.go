package main

import (
	"GoPractice/src/stringutil"
	"fmt"
	"math/rand"
)

type vertex struct {
	Lat, Long float64
}

func main() {
	fmt.Println(stringutil.Reverse("!oG ,olleH"))
	fmt.Println("My favorite number is", rand.Intn(10), "--")
	fmt.Println(add(3, 4))
	fmt.Println(swap("World", "Hello "))

	var m = make(map[string]vertex)
	m["Bell Labs"] = vertex{
		40.68433, -74.39967,
	}
	fmt.Println(m["Bell Labs"])
}

func add(x int, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}
