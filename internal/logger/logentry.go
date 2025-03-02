package logger

import (
	"errors"
	"fmt"
	"time"
)

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// LogEntry represents a single log entry with various attributes.
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`          // The time when the log entry was created.
	Level     LogLevel               `json:"level"`              // The severity level of the log entry.
	Source    string                 `json:"source"`             // The source of the log entry.
	Context   string                 `json:"context,omitempty"`  // Additional context for the log entry.
	Message   string                 `json:"message"`            // The log message.
	Tags      map[string]string      `json:"tags,omitempty"`     // Optional tags associated with the log entry.
	Metadata  map[string]interface{} `json:"metadata,omitempty"` // Optional metadata associated with the log entry.
	ProcessID int                    `json:"pid,omitempty"`      // The process ID associated with the log entry.
	Hostname  string                 `json:"hostname,omitempty"` // The hostname where the log entry was created.
	Severity  int                    `json:"severity"`           // The severity level as an integer.
	TraceID   string                 `json:"trace_id,omitempty"` // Optional trace ID for tracing logs.
	Caller    string                 `json:"caller,omitempty"`   // The caller of the log entry.
}

// NewLogEntry creates a new instance of LogEntry with the current timestamp and initialized maps.
func NewLogEntry() *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}
}

// WithLevel sets the log level for the LogEntry.
func (le *LogEntry) WithLevel(level LogLevel) *LogEntry {
	le.Level = level
	return le
}

// WithSource sets the source for the LogEntry.
func (le *LogEntry) WithSource(source string) *LogEntry {
	le.Source = source
	return le
}

// WithContext sets the context for the LogEntry.
func (le *LogEntry) WithContext(context string) *LogEntry {
	le.Context = context
	return le
}

// WithMessage sets the message for the LogEntry.
func (le *LogEntry) WithMessage(message string) *LogEntry {
	le.Message = message
	return le
}

// WithProcessID sets the process ID for the LogEntry.
func (le *LogEntry) WithProcessID(pid int) *LogEntry {
	le.ProcessID = pid
	return le
}

// WithHostname sets the hostname for the LogEntry.
func (le *LogEntry) WithHostname(hostname string) *LogEntry {
	le.Hostname = hostname
	return le
}

// WithSeverity sets the severity level for the LogEntry.
func (le *LogEntry) WithSeverity(severity int) *LogEntry {
	le.Severity = severity
	return le
}

// WithTraceID sets the trace ID for the LogEntry.
func (le *LogEntry) WithTraceID(traceID string) *LogEntry {
	le.TraceID = traceID
	return le
}

// AddTag adds a tag to the LogEntry.
func (le *LogEntry) AddTag(key, value string) *LogEntry {
	if le.Tags == nil {
		le.Tags = make(map[string]string)
	}
	le.Tags[key] = value
	return le
}

// AddMetadata adds metadata to the LogEntry.
func (le *LogEntry) AddMetadata(key string, value interface{}) *LogEntry {
	if le.Metadata == nil {
		le.Metadata = make(map[string]interface{})
	}
	le.Metadata[key] = value
	return le
}

// Validate checks if the LogEntry has all required fields set.
func (le *LogEntry) Validate() error {
	if le.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if le.Level == "" {
		return errors.New("level is required")
	}
	if le.Message == "" {
		return errors.New("message is required")
	}
	if le.Severity <= 0 {
		return errors.New("severity must be greater than zero")
	}
	return nil
}

// String returns a string representation of the LogEntry.
func (le *LogEntry) String() string {
	return fmt.Sprintf("[%s] %s - %s",
		le.Timestamp.Format(time.RFC3339),
		le.Level,
		le.Message)
}
