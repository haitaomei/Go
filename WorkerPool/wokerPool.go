package main

import (
	"fmt"
	"sync"
	"time"
)

var taskQueue chan int
var resultQueue chan int
var wg sync.WaitGroup

func main() {

	initThreadPool()

	// submit jobs
	for j := 1; j <= 5; j++ {
		wg.Add(1)
		taskQueue <- j
	}

	wg.Wait()

	for a := 1; a <= 5; a++ {
		tmp := <-resultQueue
		fmt.Println("res", tmp)
	}

	close(taskQueue)
	close(resultQueue)
}

func initThreadPool() {
	taskQueue = make(chan int, 100)
	resultQueue = make(chan int, 100)
	for w := 1; w <= 2; w++ {
		go threadPool(w, taskQueue, resultQueue)
	}
}

func threadPool(id int, taskQueue <-chan int, resultQueue chan<- int) {
	for j := range taskQueue {
		fmt.Println("worker", id, "get task", j)
		time.Sleep(100 * time.Millisecond)
		fmt.Println("worker", id, "finished task", j)
		resultQueue <- j * 2
		wg.Done()
	}
}
