package worker

import (
	"context"
	"fmt"

	"github.com/Chukwuemeka-James/go-concurrent-web-scraper/fetcher"
)

// Start launches a worker goroutine that continuously listens for jobs (URLs) on the 'jobs' channel.
// Each worker fetches the content of the URL using fetcher.FetchWithRetry and sends the result
// or error message into the 'results' channel.
// The worker stops gracefully when either:
//  1. The context is canceled (graceful shutdown).
//  2. The jobs channel is closed (no more URLs to process).
//
// Parameters:
//   - ctx: Context used to handle cancellation for graceful shutdown.
//   - id: Unique identifier for the worker (useful for logging/debugging).
//   - jobs: A read-only channel from which the worker receives URLs to process.
//   - results: A write-only channel where the worker sends results or error messages.
func Start(ctx context.Context, id int, jobs <-chan string, results chan<- string) {
	for {
		select {
		// Case 1: Listen for cancellation signal from context
		case <-ctx.Done():
			// If the context is canceled (e.g., user pressed Ctrl+C),
			// the worker logs a stop message and exits the loop.
			fmt.Printf("[Worker %d] Stopping\n", id)
			return

		// Case 2: Receive a job (URL) from the jobs channel
		case url, ok := <-jobs:
			// If the channel is closed (no more jobs), exit the worker
			if !ok {
				return
			}

			// Fetch the content of the URL with up to 3 retries
			body, err := fetcher.FetchWithRetry(url, 3)
			if err != nil {
				// If an error occurs (e.g., network issue or non-200 response after retries),
				// send an error message to the results channel
				results <- fmt.Sprintf("Worker %d: Error fetching %s: %v", id, url, err)
			} else {
				// On success, send the URL and the length of the fetched body
				results <- fmt.Sprintf("Worker %d: Fetched %s, length: %d", id, url, len(body))
			}
		}
	}
}
