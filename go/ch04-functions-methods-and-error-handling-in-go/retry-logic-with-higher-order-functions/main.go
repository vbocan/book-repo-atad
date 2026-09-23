package main

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

func retry(attempts int, delay time.Duration, operation func() error) error {
    var lastErr error
    for i := 0; i < attempts; i++ {
        if i > 0 {
            fmt.Printf("  Retrying (attempt %d/%d)...\n", i+1, attempts)
            time.Sleep(delay)
        }
        lastErr = operation()
        if lastErr == nil {
            return nil
        }
        fmt.Printf("  Attempt %d failed: %v\n", i+1, lastErr)
    }
    return fmt.Errorf("all %d attempts failed, last error: %w", attempts, lastErr)
}

// Simulates an unreliable operation that fails randomly
func unreliableOperation() error {
    if rand.Intn(3) == 0 { // ~33% chance of success
        return nil
    }
    return errors.New("temporary failure")
}

func main() {
    fmt.Println("Attempting unreliable operation...")
    err := retry(5, 500*time.Millisecond, unreliableOperation)
    if err != nil {
        fmt.Printf("Failed: %v\n", err)
    } else {
        fmt.Println("Success!")
    }
}
