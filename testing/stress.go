package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var count int64
	var wg sync.WaitGroup

	// Channel to signal workers to stop
	done := make(chan struct{})

	// Create a HTTP client (you can customize timeouts if needed)
	client := &http.Client{}

	// Worker function which keeps sending GET requests until signaled to stop
	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				resp, err := client.Get("http://localhost:9812/")
				if err == nil && resp != nil {
					// Ensure the body is closed
					resp.Body.Close()
				}
				atomic.AddInt64(&count, 1)
			}
		}
	}

	// Number of concurrent workers.
	// You can adjust this number based on your machine's capabilities.
	numWorkers := 100
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Run the test for 60 seconds.
	time.Sleep(60 * time.Second)
	close(done)
	wg.Wait()

	total := atomic.LoadInt64(&count)
	fmt.Printf("Total requests: %d\n", total)
	fmt.Printf("Requests per second: %.2f\n", float64(total)/60)
}
