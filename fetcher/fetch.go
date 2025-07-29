package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchWithRetry performs an HTTP GET request to the provided URL with a retry mechanism.
// It will attempt to fetch the URL up to 'retries' times before failing.
//
// Parameters:
//   - url: The target URL to fetch.
//   - retries: The maximum number of retry attempts if the request fails.
//
// Returns:
//   - string: The response body if the request is successful.
//   - error: An error if all retry attempts fail.
func FetchWithRetry(url string, retries int) (string, error) {
	var resp *http.Response // Will hold the HTTP response object
	var err error           // Will hold any error encountered during the request

	// Retry loop: Attempt fetching the URL up to 'retries' times
	for i := 0; i < retries; i++ {
		// Perform the HTTP GET request
		resp, err = http.Get(url)

		// If the request is successful and the server returns a 200 OK status
		if err == nil && resp.StatusCode == http.StatusOK {
			// Ensure that the response body is closed once we're done reading it
			defer resp.Body.Close()

			// Read the entire response body into memory
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				// Return an error if we fail to read the response body
				return "", fmt.Errorf("error reading body: %v", err)
			}

			// Convert the response body to a string and return it
			return string(body), nil
		}

		// If the request failed or returned a non-200 status,
		// wait before retrying (simple linear backoff: 1s, 2s, 3s...)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	// If all retries fail, return an error indicating the failure and the last error encountered
	return "", fmt.Errorf("failed after %d retries: %v", retries, err)
}
