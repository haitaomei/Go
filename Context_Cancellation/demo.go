package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func main() {
	EmittingCancellationEvent()
	TimeBasedCancellation()
	//listenCancellationEvent()
}

func listenCancellationEvent() {
	_ = http.ListenAndServe(":60080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println("processing request")
		// We use `select` to execute a piece of code depending on which channel receives a message first
		select {
		case <-time.After(2 * time.Second):
			// Simulate a message after 2 seconds meaning the request has been processed
			_, _ = w.Write([]byte("request processed"))
		case <-ctx.Done():
			// If the request gets cancelled, log it to STDERR
			fmt.Println("request cancelled")
		}
	}))
}

func EmittingCancellationEvent() {
	/*
		This demo simulate two dependent tasks.
		As soon as one failed, cancel all the others.
	*/

	// Create a new context and its cancellation function
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		err := task1(ctx)
		// If this operation returns an error cancel all operations using this context
		if err != nil {
			cancel()
		}
	}()

	// Run task2 with the same context with task1
	task2(ctx)
}

func task1(ctx context.Context) error {
	// simulate some work
	time.Sleep(100 * time.Millisecond)
	// simulate an error
	return errors.New("failed")
}

func task2(ctx context.Context) {
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("Task2 finished")
	case <-ctx.Done():
		fmt.Println("Task2 is cancelled")
	}
}

func TimeBasedCancellation() {
	// Create a new context
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 300*time.Millisecond) // With a timeout of 100 milliseconds
	// or with an absolute deadline
	//ctx, cancel := context.WithDeadline(ctx, time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC))

	// Make a request, and associate the cancellable context
	req, _ := http.NewRequest(http.MethodGet, "http://google.com", nil)
	req = req.WithContext(ctx)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	res, err := client.Do(req)
	// If the request failed, log to STDOUT
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	fmt.Println("Response received, status code:", res.StatusCode)
}
