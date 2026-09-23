package main

import (
    "fmt"
    "os"
    "strings"
    "unicode"
    "unicode/utf8"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: strstat <text>")
        os.Exit(1)
    }

    text := strings.Join(os.Args[1:], " ")

    bytes := len(text)
    runes := utf8.RuneCountInString(text)
    words := len(strings.Fields(text))

    var upper, lower, digits int
    for _, r := range text {
        switch {
        case unicode.IsUpper(r):
            upper++
        case unicode.IsLower(r):
            lower++
        case unicode.IsDigit(r):
            digits++
        }
    }

    fmt.Printf("Text:       %q\n", text)
    fmt.Printf("Bytes:      %d\n", bytes)
    fmt.Printf("Characters: %d\n", runes)
    fmt.Printf("Words:      %d\n", words)
    fmt.Printf("Uppercase:  %d\n", upper)
    fmt.Printf("Lowercase:  %d\n", lower)
    fmt.Printf("Digits:     %d\n", digits)
}
