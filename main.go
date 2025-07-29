package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Chukwuemeka-James/go-concurrent-web-scraper/worker"
)

// workerCount defines the number of concurrent workers (goroutines)
// that will process URLs from the jobs channel.
const workerCount = 5

func main() {
	// Create a buffered channel to hold jobs (URLs to scrape).
	// The buffer size of 10 allows sending jobs without blocking immediately
	// if workers are temporarily busy.
	jobs := make(chan string, 10)

	// Create a buffered channel for results from workers.
	// Workers will push either success or error messages here.
	results := make(chan string, 10)

	// Create a context with cancellation to allow graceful shutdown.
	// When "cancel" is called, all workers listening to this context will stop.
	ctx, cancel := context.WithCancel(context.Background())

	// Start a pool of worker goroutines.
	// Each worker runs the worker.Start function and listens for jobs on the "jobs" channel.
	for i := 0; i < workerCount; i++ {
		go worker.Start(ctx, i, jobs, results)
	}

	// Graceful shutdown handler (runs in a separate goroutine).
	// Listens for an OS interrupt signal (e.g. Ctrl+C) and then cancels the context.
	go func() {
		c := make(chan os.Signal, 1) // Channel to receive interrupt signals.
		signal.Notify(c, os.Interrupt)
		<-c // Wait until an interrupt signal is received.
		fmt.Println("Shutting down...")
		cancel() // Cancel the context to signal all workers to stop.
	}()

	// Feed jobs (URLs) into the jobs channel in a separate goroutine.
	// This reads URLs from "urls.txt" and sends them to workers.
	go func() {
		file, _ := os.Open("urls.txt") // Open the file containing the URLs.
		defer file.Close()

		scanner := bufio.NewScanner(file) // Create a scanner to read the file line by line.
		for scanner.Scan() {
			jobs <- scanner.Text() // Send each URL to the jobs channel.
		}
		close(jobs) // Close the jobs channel after all URLs are sent.
	}()

	// Main goroutine reads results from the results channel and prints them.
	// This loop will keep running until results are no longer being sent.
	for res := range results {
		fmt.Println(res)
	}
}
