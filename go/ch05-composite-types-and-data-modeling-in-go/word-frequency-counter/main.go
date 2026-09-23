// Command word-frequency-counter is a reference solution for Exercise 5.2.
// It reads a text file, tokenizes it into lowercase words, and prints a
// frequency table sorted by count descending, breaking ties alphabetically.
//
// Usage:
//
//	word-frequency-counter <file>
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: word-frequency-counter <file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
	defer file.Close()

	counts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := cleanWord(scanner.Text())
		if word == "" {
			continue
		}
		counts[word]++
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	printFrequencies(counts)
}

// cleanWord lowercases a whitespace-delimited token and trims leading and
// trailing punctuation, so "Hello," and "hello" count as the same word while
// internal punctuation (apostrophes, hyphens) is left alone.
func cleanWord(token string) string {
	token = strings.ToLower(token)
	return strings.TrimFunc(token, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// printFrequencies prints "word: count" lines ordered by count descending,
// with ties broken alphabetically by word.
func printFrequencies(counts map[string]int) {
	words := make([]string, 0, len(counts))
	for w := range counts {
		words = append(words, w)
	}

	sort.Slice(words, func(i, j int) bool {
		if counts[words[i]] != counts[words[j]] {
			return counts[words[i]] > counts[words[j]]
		}
		return words[i] < words[j]
	})

	for _, w := range words {
		fmt.Printf("%s: %d\n", w, counts[w])
	}
}
