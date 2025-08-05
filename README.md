## go-concurrent-web-scraper

A concurrent web scraper built in Go that demonstrates **goroutines**, **channels**, **worker pools**, **retry logic**, and **graceful shutdown**.
This project is designed to help developers move from beginner-level concurrency knowledge in Go to a **solid intermediate understanding**.

---

## Problem Statement

Web scraping is often I/O-bound and requires fetching data from multiple URLs efficiently. Fetching URLs sequentially is slow and inefficient, especially when some URLs are slow or fail intermittently.

The challenge is to:

1. **Fetch multiple URLs concurrently** using goroutines.
2. **Distribute the work** evenly across a pool of workers.
3. Handle **retry logic** for failed requests.
4. **Collect results** from multiple workers safely.
5. **Gracefully shut down** when the user interrupts the program.
6. Make the system extensible for future improvements like parsing HTML or saving data.

This project implements a **worker pool pattern with channels**, enabling concurrent URL fetching while maintaining structured and safe communication between goroutines.

---

## Project Structure

```
go-concurrent-web-scraper/
│
├── main.go                # Entry point: initializes channels, workers, and handles shutdown
├── fetcher/
│   └── fetch.go           # Handles HTTP requests with retry logic
├── worker/
│   └── pool.go            # Worker pool: reads jobs from channel and sends results
├── urls.txt               # List of URLs to scrape
└── README.md              # Project documentation
```

---

## File-by-File Explanation

### 1. `main.go`

This is the entry point of the application.

Responsibilities:

* Initializes **job** and **result** channels.
* Starts a configurable number of **workers** (goroutines) to process jobs.
* Reads URLs from `urls.txt` and feeds them into the job channel.
* Handles **graceful shutdown** using `os.Signal` and `context.WithCancel`.
* Collects results from workers and prints them to the console.

Key concepts used:

* **Channels** (`jobs`, `results`) for communication between goroutines.
* **Context cancellation** to stop workers cleanly.
* **Signal handling** (`os.Interrupt`) for graceful shutdown.

---

### 2. `fetcher/fetch.go`

This package handles **HTTP requests with retry logic**.

Responsibilities:

* Makes HTTP GET requests to a given URL.
* Retries failed requests (e.g., network errors or `5xx` responses).
* Returns the response body as a string if successful.

Key concepts used:

* `io.ReadAll` to read HTTP response bodies.
* **Exponential backoff (or linear delay)** between retries.
* Returning meaningful errors when retries are exhausted.

---

### 3. `worker/pool.go`

This package implements the **worker pool**.

Responsibilities:

* Each worker goroutine:

  * Reads URLs from the `jobs` channel.
  * Calls `fetcher.FetchWithRetry` to fetch the content.
  * Sends results (or error messages) to the `results` channel.
* Stops gracefully when `context.Context` is canceled.

Key concepts used:

* **Fan-out pattern**: multiple workers consuming from a single job channel.
* **Fan-in pattern**: results are sent into a single result channel.
* **Context cancellation** to stop workers on shutdown.

This file is read by `main.go`, and each URL is submitted as a job to the worker pool.

---

## How It Works

1. **Startup**

   * The program starts by creating two buffered channels:

     * `jobs`: holds URLs to scrape.
     * `results`: holds fetch results or errors.
   * A `context.WithCancel` is created for graceful shutdown.

2. **Worker Pool Initialization**

   * `N` workers (e.g. 5) are started using `worker.Start`.
   * Each worker listens for jobs on the `jobs` channel.

3. **Job Feeding**

   * The program reads each URL from `urls.txt` and pushes it into the `jobs` channel.
   * When all URLs are queued, the `jobs` channel is closed.

4. **Processing**

   * Each worker fetches URLs concurrently.
   * Workers call `fetcher.FetchWithRetry`, which retries failed requests (e.g. 503 or timeout).
   * Each successful or failed result is pushed to the `results` channel.

5. **Graceful Shutdown**

   * The program listens for `os.Interrupt` (CTRL+C).
   * On interrupt, it calls `cancel()` which signals all workers to stop.

6. **Result Collection**

   * The main goroutine continuously listens on `results` and prints them until all workers finish.

---

## Running the Project

1. Clone the repository:

   ```bash
   git clone https://github.com/Chukwuemeka-James/go-concurrent-web-scraper.git
   cd go-concurrent-web-scraper
   ```

2. Initialize the Go module:

   ```bash
   go mod tidy
   ```

3. Run the scraper:

   ```bash
   go run main.go
   ```

4. Outputs should look like this depending on you urls:

   ```
   Worker 1: Fetched https://example.com, length: 1256
   Worker 2: Error fetching https://thisurldoesnotexist.tld: failed after 3 retries: ...
   Worker 3: Fetched https://golang.org, length: 8043
   ...
   

---

## What You’ll Learn

By working through this project, you'll gain hands-on experience with:

* **Goroutines**: running concurrent tasks efficiently.
* **Channels**: communication between concurrent tasks.
* **Worker pools**: limiting concurrency for controlled resource usage.
* **Retry logic**: handling transient errors.
* **Context cancellation**: stopping goroutines gracefully.
* **Signal handling**: responding to user interrupts.

---

## Next Steps

You can enhance the scraper by:

* Adding timeouts for HTTP requests.
* Implementing exponential backoff for retries.
* Parsing and extracting specific HTML or JSON content.
* Logging results to a file or database.
* Adding rate-limiting to avoid overwhelming servers.
* Building a simple web UI to monitor progress.
