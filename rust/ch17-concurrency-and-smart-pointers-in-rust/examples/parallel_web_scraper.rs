//! Exercise 17.1 -- Parallel Web Scraper
//!
//! Fetches a batch of URLs concurrently: a fixed pool of worker threads pulls
//! URLs off a shared `Arc<Mutex<Vec<String>>>` queue, "fetches" each one, and
//! reports the outcome back to `main` over an `mpsc` channel. `main` collects
//! every result, joins the worker threads, and prints a summary showing that
//! the whole batch finished in noticeably less time than fetching the URLs
//! one by one would have taken.
//!
//! The fetch itself is simulated: `fetch_url` below does a `thread::sleep`
//! proportional to the URL's length instead of making a real network call,
//! so the example needs nothing beyond the standard library. See the
//! Challenge in the book for what changes if `fetch_url` performs a real
//! blocking HTTP request instead.
//!
//! Run with `cargo run --example parallel_web_scraper`.

use std::sync::{mpsc, Arc, Mutex};
use std::thread;
use std::time::{Duration, Instant};

/// How many worker threads pull from the queue at once.
const WORKER_COUNT: usize = 4;

/// One completed (simulated) fetch, as reported back to `main`.
#[derive(Debug)]
struct FetchResult {
    url: String,
    worker_id: usize,
    status: u16,
    bytes: usize,
    latency: Duration,
}

/// Stands in for a real HTTP GET. Instead of hitting the network, it sleeps
/// for a duration proportional to the URL's length -- a longer URL takes
/// longer to "download" -- and fabricates a status code and a response
/// size. Swapping this for a real `reqwest`/`ureq` call would not require
/// changing anything else in the program: the worker below only cares that
/// `fetch_url` blocks for a while and then returns a status and a size.
fn fetch_url(url: &str) -> (u16, usize) {
    let simulated_latency = Duration::from_millis(url.len() as u64 * 10);
    thread::sleep(simulated_latency);

    let status = 200;
    let bytes = 512 + url.len() * 47; // Stand-in for a response body size.
    (status, bytes)
}

/// Runs on a worker thread: repeatedly pops a URL from the shared queue,
/// fetches it, and sends the outcome back to `main` over the channel. The
/// loop ends on its own once the queue is empty -- no shutdown signal is
/// needed because `Vec::pop` returning `None` already says "no work left".
fn worker(id: usize, queue: Arc<Mutex<Vec<String>>>, tx: mpsc::Sender<FetchResult>) {
    loop {
        // Hold the lock only long enough to pop one URL, so the other
        // workers are never blocked waiting on this one's network "call".
        let next = {
            let mut q = queue.lock().unwrap();
            q.pop()
        };

        let url = match next {
            Some(url) => url,
            None => break, // Queue is empty; this worker is done.
        };

        let started = Instant::now();
        let (status, bytes) = fetch_url(&url);
        let latency = started.elapsed();

        tx.send(FetchResult { url, worker_id: id, status, bytes, latency })
            .expect("main thread dropped the receiver before the run finished");
    }
}

fn main() {
    let urls: Vec<String> = vec![
        "http://a.io",
        "http://example.com",
        "http://example.com/health",
        "http://example.com/api/v1/users",
        "http://example.com/api/v1/users/12345",
        "http://example.com/api/v1/users/12345/profile",
        "http://example.com/static/logo.png",
        "http://example.com/static/images/banner-large.jpg",
        "http://news.example.org/latest",
        "http://news.example.org/latest/breaking-news-today",
        "http://example.com/blog/2024/09/rust-concurrency",
        "http://example.com/blog/2024/09/rust-concurrency-patterns-explained",
        "http://example.com/search?q=rust+programming",
        "http://example.com/search?q=rust+programming+language+tutorial",
        "http://example.com/docs/reference/api/v3/endpoints",
        "http://example.com/docs/reference/api/v3/endpoints/list/all/items",
    ]
    .into_iter()
    .map(String::from)
    .collect();

    let url_count = urls.len();

    // What fetching every URL one after another, on a single thread, would
    // have cost. This is the number the concurrent run has to beat.
    let serial_equivalent: Duration = urls
        .iter()
        .map(|url| Duration::from_millis(url.len() as u64 * 10))
        .sum();

    // The shared work queue: every worker thread sees the same `Vec`
    // through its own `Arc` handle, and the `Mutex` makes popping from it
    // safe even though several threads do it at once.
    let queue = Arc::new(Mutex::new(urls));

    // The channel workers use to report results back to `main`.
    let (tx, rx) = mpsc::channel::<FetchResult>();

    println!(
        "Fetching {} URLs with {} worker threads...\n",
        url_count, WORKER_COUNT
    );

    let started = Instant::now();

    let mut handles = Vec::with_capacity(WORKER_COUNT);
    for id in 0..WORKER_COUNT {
        let queue = Arc::clone(&queue);
        let tx = tx.clone();
        handles.push(thread::spawn(move || worker(id, queue, tx)));
    }

    // Drop the original sender. Each worker holds its own clone, and once
    // every clone is gone -- i.e. every worker has finished -- the
    // receiver's iterator below ends on its own.
    drop(tx);

    // Collect every result. This blocks until the channel closes, which
    // happens naturally when the last worker thread's `tx` clone is dropped.
    let mut results: Vec<FetchResult> = rx.iter().collect();

    for handle in handles {
        handle.join().expect("a worker thread panicked");
    }

    let elapsed = started.elapsed();

    // Sort for a stable, readable report -- the arrival order over the
    // channel reflects thread scheduling, not anything meaningful.
    results.sort_by(|a, b| a.url.len().cmp(&b.url.len()));

    for r in &results {
        println!(
            "  worker {}  [{}]  {:>5} bytes  {:>4} ms  {}",
            r.worker_id,
            r.status,
            r.bytes,
            r.latency.as_millis(),
            r.url,
        );
    }

    println!();
    println!(
        "Fetched {} URLs in {:?} (concurrent, {} workers)",
        url_count, elapsed, WORKER_COUNT
    );
    println!("Serial equivalent would have taken {:?}", serial_equivalent);
    println!(
        "Speedup: {:.2}x",
        serial_equivalent.as_secs_f64() / elapsed.as_secs_f64()
    );
}
