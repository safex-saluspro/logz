package logger

import (
	"errors"
	"fmt"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Source    string                 `json:"source"`
	Context   string                 `json:"context,omitempty"`
	Message   string                 `json:"message"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ProcessID int                    `json:"pid,omitempty"`
	Hostname  string                 `json:"hostname,omitempty"`
	Severity  int                    `json:"severity"`
	TraceID   string                 `json:"trace_id,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
}

func NewLogEntry() *LogEntry {
	// Obter o nome da função que chamou a função NewLogEntry.
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}
}

func (le *LogEntry) WithLevel(level LogLevel) *LogEntry {
	le.Level = level
	return le
}
func (le *LogEntry) WithSource(source string) *LogEntry {
	le.Source = source
	return le
}
func (le *LogEntry) WithContext(context string) *LogEntry {
	le.Context = context
	return le
}
func (le *LogEntry) WithMessage(message string) *LogEntry {
	le.Message = message
	return le
}
func (le *LogEntry) WithProcessID(pid int) *LogEntry {
	le.ProcessID = pid
	return le
}
func (le *LogEntry) WithHostname(hostname string) *LogEntry {
	le.Hostname = hostname
	return le
}
func (le *LogEntry) WithSeverity(severity int) *LogEntry {
	le.Severity = severity
	return le
}
func (le *LogEntry) WithTraceID(traceID string) *LogEntry {
	le.TraceID = traceID
	return le
}
func (le *LogEntry) AddTag(key, value string) *LogEntry {
	if le.Tags == nil {
		le.Tags = make(map[string]string)
	}
	le.Tags[key] = value
	return le
}
func (le *LogEntry) AddMetadata(key string, value interface{}) *LogEntry {
	if le.Metadata == nil {
		le.Metadata = make(map[string]interface{})
	}
	le.Metadata[key] = value
	return le
}
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
func (le *LogEntry) String() string {
	return fmt.Sprintf("[%s] %s - %s",
		le.Timestamp.Format(time.RFC3339),
		le.Level,
		le.Message)
}
