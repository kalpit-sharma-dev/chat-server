package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
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
func worker(id int, jobs <-chan Request, wg *sync.WaitGroup, logger *log.Logger) {
	defer wg.Done()
	for req := range jobs {
		logger.Printf("Worker %d: Processing request to %s\n", id, req.URL)

		// Create HTTP POST request
		request, err := http.NewRequest("POST", req.URL, bytes.NewBuffer(req.Payload))
		if err != nil {
			logger.Printf("Worker %d: Error creating request: %v\n", id, err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := req.Client.Do(request)
		if err != nil {
			logger.Printf("Worker %d: Error making request: %v\n", id, err)
			continue
		}
		logger.Printf("Worker %d: Received response with status: %s\n", id, resp.Status)
		resp.Body.Close()
	}
}

func main() {
	const numWorkers = 5
	const numJobs = 10

	// Create log file
	logFile, err := os.OpenFile("worker.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error creating log file:", err)
		return
	}
	defer logFile.Close()

	// Initialize logger
	logger := log.New(logFile, "", log.LstdFlags)

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
		go worker(i, jobs, &wg, logger)
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
	logger.Println("All requests completed.")
}
