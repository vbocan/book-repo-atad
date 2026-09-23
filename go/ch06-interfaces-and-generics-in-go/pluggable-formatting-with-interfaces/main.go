// Exercise 6.3: Pluggable Formatting with Interfaces
//
// This logging system accepts different formatters through a single
// Formatter interface. Logger holds a Formatter and delegates all
// rendering to it, so the same Log calls produce different output
// depending on which formatter was injected at construction. Adding a
// new output format (CSV, XML, syslog, ...) only requires a new
// Formatter implementation; Logger and its callers never change.
//
// Chapter 15 reimplements this same exercise in Rust using traits, so
// the shapes here (a Formatter contract, a handful of implementations,
// and a Logger that composes with one of them) are kept deliberately
// simple to carry over cleanly.
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry is a single event passed to a Formatter for rendering.
type LogEntry struct {
	Level   string
	Message string
	Time    time.Time
}

// Formatter renders a LogEntry as a string.
type Formatter interface {
	Format(entry LogEntry) string
}

// PlainFormatter renders entries as simple, greppable text lines.
type PlainFormatter struct{}

func (PlainFormatter) Format(entry LogEntry) string {
	return fmt.Sprintf("[%s] %s: %s",
		entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message)
}

// JSONFormatter renders entries as single-line JSON objects, suitable
// for log aggregation systems that expect structured logs.
type JSONFormatter struct{}

func (JSONFormatter) Format(entry LogEntry) string {
	data, err := json.Marshal(struct {
		Level   string    `json:"level"`
		Message string    `json:"message"`
		Time    time.Time `json:"time"`
	}{entry.Level, entry.Message, entry.Time})
	if err != nil {
		return fmt.Sprintf(`{"level":"ERROR","message":"failed to format log entry: %s"}`, err)
	}
	return string(data)
}

// ANSI escape codes for terminal colors. No third-party color library is
// needed since the sequences are just bytes the terminal interprets.
const (
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiGreen  = "\033[32m"
	ansiReset  = "\033[0m"
)

// ColorFormatter renders entries as plain text, with the level word
// colorized for terminals that understand ANSI escape codes: red for
// ERROR, yellow for WARN, green for INFO.
type ColorFormatter struct{}

func (ColorFormatter) Format(entry LogEntry) string {
	color := ansiReset
	switch entry.Level {
	case "ERROR":
		color = ansiRed
	case "WARN":
		color = ansiYellow
	case "INFO":
		color = ansiGreen
	}
	return fmt.Sprintf("[%s] %s%s%s: %s",
		entry.Time.Format("2006-01-02 15:04:05"), color, entry.Level, ansiReset, entry.Message)
}

// Logger records events and delegates all rendering to a Formatter.
type Logger struct {
	formatter Formatter
}

// NewLogger creates a Logger that renders entries with formatter.
func NewLogger(formatter Formatter) *Logger {
	return &Logger{formatter: formatter}
}

// Log records one entry at the given level and prints it through the
// Logger's Formatter.
func (l *Logger) Log(level, message string) {
	entry := LogEntry{Level: level, Message: message, Time: time.Now()}
	fmt.Println(l.formatter.Format(entry))
}

func main() {
	formatters := []struct {
		name      string
		formatter Formatter
	}{
		{"Plain", PlainFormatter{}},
		{"JSON", JSONFormatter{}},
		{"Color", ColorFormatter{}},
	}

	for _, f := range formatters {
		fmt.Printf("--- %s formatter ---\n", f.name)
		logger := NewLogger(f.formatter)
		logger.Log("INFO", "service started")
		logger.Log("WARN", "cache miss rate above threshold")
		logger.Log("ERROR", "database connection lost")
	}
}
