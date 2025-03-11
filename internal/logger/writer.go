package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// LogFormatter defines the contract for formatting log entries.
type LogFormatter interface {
	Format(entry LogzEntry) (string, error)
}

// JSONFormatter formats the log in JSON format.
type JSONFormatter struct{}

// Format converts the log entry to JSON.
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
func (f *TextFormatter) Format(entry LogzEntry) (string, error) {
	noColor := os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows"
	icon, levelStr := "", ""

	// Configure cores e ícones por nível
	if !noColor {
		switch entry.GetLevel() {
		case DEBUG:
			icon, levelStr = "\033[34m🐛\033[0m", "\033[34mDEBUG\033[0m"
		case INFO:
			icon, levelStr = "\033[32mℹ️\033[0m", "\033[32mINFO\033[0m"
		case WARN:
			icon, levelStr = "\033[33m⚠️\033[0m", "\033[33mWARN\033[0m"
		case ERROR:
			icon, levelStr = "\033[31m❌\033[0m", "\033[31mERROR\033[0m"
		case FATAL:
			icon, levelStr = "\033[35m💀\033[0m", "\033[35mFATAL\033[0m"
		default:
			icon, levelStr = string(entry.GetLevel()), ""
		}
	} else {
		icon, levelStr = string(entry.GetLevel()), ""
	}

	// Header compacto
	header := fmt.Sprintf("[%s] %s %s", entry.GetTimestamp().Format(time.RFC3339), icon, levelStr)

	// Contexto e Metadados
	metadata := ""
	if len(entry.GetMetadata()) > 0 {
		if entry.GetLevel() == DEBUG || entry.GetMetadata()["showContext"] == true {
			metadata = fmt.Sprintf("\n  %s", formatMetadata(entry.GetMetadata()))
		} /* else {
			metadata = fmt.Sprintf(" -> %s", formatMetadata(entry.GetMetadata()))
		}*/
	}

	// Formato final
	return fmt.Sprintf("%s | %s%s", header, entry.GetMessage(), metadata), nil
}

// LogWriter defines the contract for writing logs.
type LogWriter interface {
	Write(entry LogzEntry) error
}

// DefaultWriter implements LogWriter using an io.Writer and a LogFormatter.
type DefaultWriter struct {
	out       io.Writer
	formatter LogFormatter
}

// NewDefaultWriter creates a new instance of DefaultWriter.
func NewDefaultWriter(out io.Writer, formatter LogFormatter) *DefaultWriter {
	return &DefaultWriter{
		out:       out,
		formatter: formatter,
	}
}

// Write formats the entry and writes it to the configured destination.
func (w *DefaultWriter) Write(entry LogzEntry) error {
	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w.out, formatted)
	return err
}

func formatMetadata(metadata map[string]interface{}) string {
	if len(metadata) == 0 {
		return ""
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return ""
	}
	return string(data)
}
