package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Request struct to hold API details
type Request struct {
	Client  *http.Client
	URL     string
	Payload []byte
}

// Worker function to process requests
func worker(id int, jobs <-chan Request, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range jobs {
		// Create HTTP POST request
		request, err := http.NewRequest("POST", req.URL, bytes.NewBuffer(req.Payload))
		if err != nil {
			fmt.Printf("Worker %d: Error creating request: %v\n", id, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := req.Client.Do(request)
		if err != nil {
			fmt.Printf("Worker %d: Error making request: %v\n", id, err)
			continue
		}
		fmt.Printf("Worker %d: Received response with status: %s\n", id, resp.Status)
		resp.Body.Close()
	}
}

func main() {
	const numWorkers = 5
	const numJobs = 10

	// Create HTTP client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Channel for worker pool jobs
	jobs := make(chan Request, numJobs)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// API Endpoint
	apiURL := "http://localhost:8080/api/users"
	payload := []byte(`{"name":"John Doe","age":30}`)

	// Send jobs to workers
	for i := 0; i < numJobs; i++ {
		jobs <- Request{Client: client, URL: apiURL, Payload: payload}
	}

	// Close jobs channel after all jobs are sent
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All requests completed.")
}
