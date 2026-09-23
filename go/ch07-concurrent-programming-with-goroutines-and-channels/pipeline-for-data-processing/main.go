// Command pipeline-for-data-processing is the reference solution for
// Exercise 7.2. It builds a four-stage pipeline out of goroutines connected
// by channels: generate the integers from 2 to 100, keep only the primes,
// square what is left, and sum the squares. The shape is exactly the one
// developed in Section 7.7.2 ("Pipeline Pattern") — generate, filter, and
// square are reused almost verbatim from that section — extended here with
// a terminal sum stage that reduces the pipeline to a single value.
//
// For the range 2..100 the expected result is:
//
//	2² + 3² + 5² + 7² + 11² + ... + 97² = 65796
//
// Run it with:
//
//	go run main.go
package main

import "fmt"

// generate emits every integer in [start, end] on its own goroutine and
// closes the channel when it is done. Identical in shape to `generate` in
// Section 7.7.2.
func generate(start, end int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := start; n <= end; n++ {
			out <- n
		}
	}()
	return out
}

// filter reads from in and forwards only the values that satisfy predicate.
// Same generic shape as Section 7.7.2's `filter`; this program instantiates
// it with isPrime instead of the evenness check used there.
func filter(in <-chan int, predicate func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if predicate(n) {
				out <- n
			}
		}
	}()
	return out
}

// square reads from in and forwards each value multiplied by itself, as in
// Section 7.7.2.
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// sum is the pipeline's terminal stage. It drains in completely, adds every
// value it receives, and sends the single running total on the returned
// channel before closing it. Doing the reduction in its own goroutine,
// rather than ranging over the upstream channel directly in main, keeps
// every stage of the pipeline — including this last one — the same shape:
// a goroutine that owns one channel and closes it when its work is done.
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		total := 0
		for n := range in {
			total += n
		}
		out <- total
	}()
	return out
}

// isPrime reports whether n is a prime number, using trial division up to
// sqrt(n). The pipeline only ever evaluates this against small values
// (2..100), so a sieve would be overkill.
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func main() {
	const start, end = 2, 100

	// generate -> filter(isPrime) -> square -> sum: four goroutines, each
	// wired to the next by a channel, running concurrently as values flow
	// from left to right.
	numbers := generate(start, end)
	primes := filter(numbers, isPrime)
	squares := square(primes)
	total := <-sum(squares)

	fmt.Printf("Sum: %d\n", total)
}
