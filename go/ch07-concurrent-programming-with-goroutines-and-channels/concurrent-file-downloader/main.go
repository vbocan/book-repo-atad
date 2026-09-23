// Command concurrent-file-downloader is the reference solution for
// Exercise 7.1. It downloads a list of URLs concurrently through a worker
// pool, reports progress for each completed download over a channel, and
// wraps the whole run in a context.Context that can be canceled either by
// an overall deadline or by an interrupt signal (Ctrl+C) — the worker pool
// pattern from Section 7.7.3 combined with the cancellation machinery from
// Section 7.6.
//
// Usage:
//
//	go run main.go [workers] [url ...]
//
// With no arguments it downloads a small default set of fast, reliable
// public URLs using 3 workers. An optional leading numeric argument sets
// the worker pool size; any remaining arguments replace the default URL
// list.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
)

// defaultURLs is used when the program is invoked with no URL arguments.
// These are small, fast, and reliably reachable, which makes them a good
// smoke test for the download pipeline itself rather than for any one
// site's content.
var defaultURLs = []string{
	"https://go.dev",
	"https://example.com",
	"https://www.rust-lang.org",
	"https://httpbin.org/get",
}

const (
	defaultWorkers = 3 // worker pool size when none is given on the command line

	perRequestTimeout = 10 * time.Second // budget for a single download
	overallTimeout    = 30 * time.Second // budget for the whole run
)

// job is one unit of work handed to a worker: a single URL to fetch, with
// its position in the original list so progress lines can be matched back
// to the request that produced them.
type job struct {
	id  int
	url string
}

// downloadResult is what a worker reports back after attempting a
// download. When err is nil, status and bytes describe the response;
// when err is non-nil, the request failed or was canceled and the other
// fields are meaningless.
type downloadResult struct {
	id     int
	url    string
	status string
	bytes  int64
	err    error
	dur    time.Duration
}

func main() {
	workers, urls := parseArgs(os.Args[1:])

	// The root context for the whole run carries an overall deadline, so
	// the program cannot hang forever even if something never responds,
	// and it is also canceled early if the user sends an interrupt. Either
	// signal reaches every in-flight download through ctx.Done(), because
	// each request below is tied to a context derived from ctx.
	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopSignal()
	ctx, cancel := context.WithTimeout(ctx, overallTimeout)
	defer cancel()

	exitCode := run(ctx, workers, urls, os.Stdout)
	os.Exit(exitCode)
}

// run drives the worker pool end to end: it starts the workers, feeds them
// jobs, prints one line per completed download as progress reports arrive,
// and prints a final summary. It returns a process exit code so that main
// stays a thin wrapper around context and signal setup.
func run(ctx context.Context, workers int, urls []string, out io.Writer) int {
	jobs := make(chan job)
	results := make(chan downloadResult)

	// Fan-out: start the worker pool. Each worker pulls jobs until the
	// jobs channel is closed or the context is canceled, whichever comes
	// first.
	var wg sync.WaitGroup
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, id, jobs, results)
		}(w)
	}

	// Feed the jobs channel from its own goroutine so that sending jobs
	// can never deadlock against the unbuffered results channel below:
	// both directions make progress independently.
	go func() {
		defer close(jobs)
		for i, u := range urls {
			select {
			case jobs <- job{id: i + 1, url: u}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Fan-in: close results once every worker has returned, so the range
	// loop below terminates instead of blocking forever.
	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Fprintf(out, "Downloading %d URL(s) with %d worker(s) (overall budget %s)...\n\n",
		len(urls), workers, overallTimeout)

	var (
		succeeded  int
		failed     int
		totalBytes int64
		start      = time.Now()
	)

	// Progress reporting: print a line the moment each result arrives,
	// rather than waiting for the whole batch, so the caller can see
	// downloads complete in whatever order the workers finish them.
	for r := range results {
		if r.err != nil {
			failed++
			fmt.Fprintf(out, "[FAILED] #%d %-28s after %-8s: %v\n",
				r.id, r.url, r.dur.Round(time.Millisecond), r.err)
			continue
		}
		succeeded++
		totalBytes += r.bytes
		fmt.Fprintf(out, "[ok]     #%d %-28s %s, %d bytes in %s\n",
			r.id, r.url, r.status, r.bytes, r.dur.Round(time.Millisecond))
	}

	fmt.Fprintf(out, "\nSummary: %d succeeded, %d failed, %d bytes total, wall time %s\n",
		succeeded, failed, totalBytes, time.Since(start).Round(time.Millisecond))

	if err := ctx.Err(); err != nil {
		fmt.Fprintf(out, "Run ended via context: %v\n", err)
	}

	if failed > 0 {
		return 1
	}
	return 0
}

// worker pulls jobs from jobs, downloads each URL, and sends one
// downloadResult per job on results, until jobs is closed or ctx is
// canceled. The send to results is a plain, unconditional send: main
// ranges over results until every worker has returned and the channel is
// closed behind them (see the fan-in goroutine in run), so a worker can
// never be left blocked here with nothing to receive it, and every job a
// worker actually started — including one that ends in cancellation — is
// guaranteed to be reported exactly once.
func worker(ctx context.Context, id int, jobs <-chan job, results chan<- downloadResult) {
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				return
			}
			results <- download(ctx, j)
		case <-ctx.Done():
			return
		}
	}
}

// download performs a single HTTP GET under a per-request timeout derived
// from ctx. http.Get has no way to accept a context, so we build the
// request with http.NewRequestWithContext instead; that is what lets a
// canceled or expired context abort an in-flight request instead of
// merely being checked before the call starts. Deriving reqCtx from ctx
// also means an interrupt or overall-deadline cancellation of the parent
// takes effect immediately, even mid-download.
func download(ctx context.Context, j job) downloadResult {
	start := time.Now()

	reqCtx, cancel := context.WithTimeout(ctx, perRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, j.url, nil)
	if err != nil {
		return downloadResult{id: j.id, url: j.url, err: err, dur: time.Since(start)}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return downloadResult{id: j.id, url: j.url, err: err, dur: time.Since(start)}
	}
	defer resp.Body.Close()

	// Read (and discard) the full body so bytes reflects what was actually
	// transferred, not just the headers, and so the connection can be
	// reused. A slow body is still bounded by reqCtx: Copy returns with an
	// error as soon as the context expires.
	n, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return downloadResult{id: j.id, url: j.url, err: err, dur: time.Since(start)}
	}

	return downloadResult{
		id:     j.id,
		url:    j.url,
		status: resp.Status,
		bytes:  n,
		dur:    time.Since(start),
	}
}

// parseArgs interprets os.Args[1:] as an optional leading worker count
// followed by zero or more URLs: `[workers] [url ...]`. With no arguments
// it returns the default worker count and the default URL list; with a
// non-numeric first argument, that argument is treated as a URL and the
// default worker count is used.
func parseArgs(args []string) (workers int, urls []string) {
	workers = defaultWorkers
	if len(args) > 0 {
		if n, err := strconv.Atoi(args[0]); err == nil {
			workers = n
			args = args[1:]
		}
	}
	if workers < 1 {
		workers = 1
	}
	if len(args) == 0 {
		return workers, defaultURLs
	}
	return workers, args
}
