package logger

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"os"
	"runtime"
)

// LogFormatter defines the contract for formatting log entries.
type LogFormatter interface {
	// Format converts a log entry to a formatted string.
	// Returns the formatted string and an error if formatting fails.
	Format(entry LogzEntry) (string, error)
}

// JSONFormatter formats the log in JSON format.
type JSONFormatter struct{}

// Format converts the log entry to JSON.
// Returns the JSON string and an error if marshalling fails.
func (f *JSONFormatter) Format(entry LogzEntry) (string, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TextFormatter formats the log in plain text.
type TextFormatter struct{}

// Format converts the log entry to a formatted string with colors and icons.
// Returns the formatted string and an error if formatting fails.
func (f *TextFormatter) Format(entry LogzEntry) (string, error) {

	// Check for environment variables
	noColor := os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows"
	noIcon := os.Getenv("LOGZ_NO_ICON") != ""

	icon, levelStr := "", ""

	if !noIcon {
		switch entry.GetLevel() {
		case DEBUG:
			icon = "\033[34m🐛\033[0m "
		case INFO:
			icon = "\033[32mℹ️\033[0m "
		case WARN:
			icon = "\033[33m⚠️\033[0m "
		case ERROR:
			icon = "\033[31m❌\033[0m "
		case FATAL:
			icon = "\033[35m💀\033[0m "
		default:
			icon = ""
		}
	} else {
		icon = ""
	}

	// Configure colors and icons by level
	if !noColor {
		switch entry.GetLevel() {
		case DEBUG:
			levelStr = "\033[34mDEBUG\033[0m"
		case INFO:
			levelStr = "\033[32mINFO\033[0m"
		case WARN:
			levelStr = "\033[33mWARN\033[0m"
		case ERROR:
			levelStr = "\033[31mERROR\033[0m"
		case FATAL:
			levelStr = "\033[35mFATAL\033[0m"
		default:
			levelStr = string(entry.GetLevel())
		}
	} else {
		levelStr = string(entry.GetLevel())
	}

	systemLocale := os.Getenv("LANG")
	tag, _ := language.Parse(systemLocale)
	p := message.NewPrinter(tag)

	withTimestamp := false

	// Context and Metadata
	metadata := ""
	if len(entry.GetMetadata()) > 0 {
		if entry.GetLevel() == DEBUG || entry.GetMetadata()["showContext"] == "true" {
			metadata = fmt.Sprintf("\n%s", formatMetadata(entry))
		}
		if entry.GetMetadata()["showTimestamp"] == "true" {
			withTimestamp = true
		}
	}

	// Determine if timestamp should be included
	willTimeStamp := os.Getenv("LOGZ_TIMESTAMP") == "true" || withTimestamp
	timestamp := ""
	if willTimeStamp {
		timestamp = fmt.Sprintf("[%s]", entry.GetTimestamp().Format(p.Sprintf("%d-%m-%Y %H:%M:%S")))
	}

	// Construct the header
	header := fmt.Sprintf("%s [%s] %s - ", timestamp, levelStr, icon)

	// Return the formatted log entry
	return fmt.Sprintf("%s%s%s", header, entry.GetMessage(), metadata), nil
}

// LogWriter defines the contract for writing logs.
type LogWriter interface {
	// Write writes a formatted log entry.
	// Returns an error if writing fails.
	Write(entry LogzEntry) error
}

// DefaultWriter implements LogWriter using an io.Writer and a LogFormatter.
type DefaultWriter struct {
	out       io.Writer
	formatter LogFormatter
}

// NewDefaultWriter creates a new instance of DefaultWriter.
// Takes an io.Writer and a LogFormatter as parameters.
func NewDefaultWriter(out io.Writer, formatter LogFormatter) *DefaultWriter {
	return &DefaultWriter{
		out:       out,
		formatter: formatter,
	}
}

// Write formats the entry and writes it to the configured destination.
// Returns an error if formatting or writing fails.
func (w *DefaultWriter) Write(entry LogzEntry) error {
	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w.out, formatted)
	return err
}

// formatMetadata converts metadata to a JSON string.
// Returns the JSON string or an empty string if marshalling fails.
func formatMetadata(entry LogzEntry) string {
	metadata := entry.GetMetadata()
	if len(metadata) == 0 {
		return ""
	}
	prefix := "Context:\n"
	for k, v := range metadata {
		if k == "showContext" {
			continue
		}
		prefix += fmt.Sprintf("  - %s: %v\n", k, v)
	}
	return prefix
}
