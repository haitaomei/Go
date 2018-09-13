package main

import (
	"fmt"
	"time"
)

func workerPool(id int, taskQueue <-chan int, resultQueue chan<- int) {
	for j := range taskQueue {
		fmt.Println("worker", id, "get task", j)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("worker", id, "finished task", j)
		resultQueue <- j * 2
	}
}

func main() {
	taskQueue := make(chan int, 100)
	resultQueue := make(chan int, 100)
	for w := 1; w <= 2; w++ {
		go workerPool(w, taskQueue, resultQueue)
	}

	for j := 1; j <= 5; j++ {
		taskQueue <- j
	}

	time.Sleep(3 * time.Second)

	for j := 1; j <= 3; j++ {
		taskQueue <- j + 5
	}
	close(taskQueue)

	for a := 1; a <= 8; a++ {
		tmp := <-resultQueue
		fmt.Println("res", tmp)
	}
}
