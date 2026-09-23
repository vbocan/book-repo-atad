package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: filestats <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", filename, err)
		os.Exit(1)
	}

	lines, words, chars := countStats(string(content))

	fmt.Printf("Lines:      %d\n", lines)
	fmt.Printf("Words:      %d\n", words)
	fmt.Printf("Characters: %d\n", chars)
}

// countStats reports the number of lines, words, and characters (runes) in
// text. A final line that is not terminated by a newline still counts.
func countStats(text string) (lines, words, chars int) {
	lines = strings.Count(text, "\n")
	if len(text) > 0 && !strings.HasSuffix(text, "\n") {
		lines++
	}
	words = len(strings.Fields(text))
	chars = utf8.RuneCountInString(text)
	return lines, words, chars
}
