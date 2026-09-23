package main

import (
	"fmt"
	"os"
	"strconv"
)

// isPrime reports whether n is a prime number.
func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: numclass <integer>")
		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid integer: %s\n", os.Args[1])
		os.Exit(1)
	}

	var sign string
	switch {
	case n > 0:
		sign = "positive"
	case n < 0:
		sign = "negative"
	default:
		sign = "zero"
	}

	parity := "even"
	if n%2 != 0 {
		parity = "odd"
	}

	result := fmt.Sprintf("%d is %s, %s", n, sign, parity)
	if n > 0 && isPrime(n) {
		result += ", prime"
	}

	fmt.Println(result)
}
