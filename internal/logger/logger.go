package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogMode string
type LogFormat string

const (
	ModeService    LogMode = "service"    // Indicates that the logger is being used by a detached process
	ModeStandalone LogMode = "standalone" // Indicates that the logger is being used locally (e.g., CLI)
)

var logLevels = map[LogLevel]int{
	DEBUG: 1,
	INFO:  2,
	WARN:  3,
	ERROR: 4,
	FATAL: 5,
}

// Logger represents a logger with configuration and metadata.
type Logger struct {
	level    LogLevel
	writer   LogWriter
	config   Config
	metadata map[string]interface{}
	mode     LogMode // Mode control: service or standalone
}

// NewLogger creates a new instance of Logger with the provided configuration.
func NewLogger(config Config) *Logger {
	// Set the log level from the Config
	level := LogLevel(config.Level()) // Method config.Level() returns the log level as a string

	var out *os.File
	if config.DefaultLogPath() == "stdout" {
		out = os.Stdout
	} else {
		// Ensure the log file exists and has the correct permissions
		if _, err := os.Stat(config.DefaultLogPath()); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(config.DefaultLogPath()), 0755); err != nil {
				log.Printf("Error creating log directory: %v\nRedirecting to stdout...\n", err)
				out = os.Stdout
			} else {
				out, err = os.OpenFile(config.DefaultLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Printf("Error opening log file: %v\nRedirecting to stdout...\n", err)
					out = os.Stdout
				}
			}
		} else {
			out, err = os.OpenFile(config.DefaultLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening log file: %v\nRedirecting to stdout...\n", err)
				out = os.Stdout
			}
		}
	}

	// Initialize the formatter (JSON or text)
	var formatter LogFormatter
	if config.Format() == "json" {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{}
	}
	writer := NewDefaultWriter(out, formatter)

	// Read the mode from Config
	mode := config.Mode()
	if mode != ModeService && mode != ModeStandalone {
		mode = ModeStandalone // Default to standalone if not specified
	}

	return &Logger{
		level:    level,
		writer:   writer,
		config:   config,
		metadata: make(map[string]interface{}),
		mode:     mode,
	}
}

// SetMetadata sets a metadata key-value pair for the Logger.
func (l *Logger) SetMetadata(key string, value interface{}) {
	l.metadata[key] = value
}

// shouldLog checks if the log level should be logged.
func (l *Logger) shouldLog(level LogLevel) bool {
	return logLevels[level] >= logLevels[l.level]
}

// log logs a message with the specified level and context.
func (l *Logger) log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	timestamp := time.Now().UTC()
	caller := getCallerInfo(3)

	entry := NewLogEntry().
		WithLevel(level).
		WithMessage(msg).
		WithSeverity(logLevels[level])
	entry.Timestamp = timestamp
	entry.Caller = caller

	// Merge global and local metadata
	finalContext := mergeContext(l.metadata, ctx)
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	// Write the log using the configured writer
	if err := l.writer.Write(entry); err != nil {
		log.Printf("Error writing log: %v", err)
	}

	// Only in service mode, notify via Notifiers
	if l.mode == ModeService && l.config != nil {
		for _, name := range l.config.NotifierManager().ListNotifiers() {
			if notifier, ok := l.config.NotifierManager().GetNotifier(name); ok {
				if notifier != nil {
					ntf := notifier
					if ntfErr := ntf.Notify(entry); ntfErr != nil {
						log.Printf("Error notifying %s: %v", name, ntfErr)
					}
				}
			}
		}
	}

	// Update metrics in PrometheusManager, if enabled
	if l.mode == ModeService {
		pm := GetPrometheusManager()
		if pm.IsEnabled() {
			pm.IncrementMetric("logs_total", 1)
			pm.IncrementMetric("logs_total_"+string(level), 1)
		}
	}

	// Terminate the process in case of FATAL log
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message with context.
func (l *Logger) Debug(msg string, ctx map[string]interface{}) { l.log(DEBUG, msg, ctx) }

// Info logs an info message with context.
func (l *Logger) Info(msg string, ctx map[string]interface{}) { l.log(INFO, msg, ctx) }

// Warn logs a warning message with context.
func (l *Logger) Warn(msg string, ctx map[string]interface{}) { l.log(WARN, msg, ctx) }

// Error logs an error message with context.
func (l *Logger) Error(msg string, ctx map[string]interface{}) { l.log(ERROR, msg, ctx) }

// Fatal logs a fatal message with context and terminates the process.
func (l *Logger) Fatal(msg string, ctx map[string]interface{}) { l.log(FATAL, msg, ctx) }

// getCallerInfo returns the caller information for the log entry.
func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", trimFilePath(file), line, funcName)
}

// trimFilePath trims the file path to show only the last two segments.
func trimFilePath(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return filePath
}

// mergeContext merges global and local context maps.
func mergeContext(global, local map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range global {
		merged[k] = v
	}
	for k, v := range local {
		merged[k] = v
	}
	return merged
}
