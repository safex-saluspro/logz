package cmd

import (
	"fmt"
	"time"
)

// LogLevel representa os níveis de log suportados
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry interface {
	Timestamp(timestamp *time.Time) time.Time
	Level(level string) string
	Source(source string) string
	Context(context string) string
	Message(message string) string
	Tags(tags map[string]string) map[string]string
	Metadata(metadata map[string]any) map[string]any
	ProcessID(pid int) int
	Hostname(hostname string) string
	Severity(severity int) int
	TraceID(traceID string) string
	AddTag(key, value string)
	Validate() error
}
type logEntryImpl struct {
	VlTimestamp *time.Time        `json:"timestamp"`
	VlLevel     LogLevel          `json:"level"`
	VlSource    string            `json:"source"`
	VlContext   string            `json:"context,omitempty"`
	VlMessage   string            `json:"message"`
	VlTags      map[string]string `json:"tags"`
	VlMetadata  map[string]any    `json:"metadata"`
	VlProcessID int               `json:"pid,omitempty"`
	VlHostname  string            `json:"hostname"`
	VlSeverity  int               `json:"severity"`
	VlTraceID   string            `json:"trace_id,omitempty"`
	Caller      string            `json:"caller,omitempty"`
}
type LogRegistry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Caller    string                 `json:"caller,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

func (l *logEntryImpl) Timestamp(timestamp *time.Time) time.Time {
	if l.VlTimestamp == nil {
		if timestamp != nil {
			l.VlTimestamp = timestamp
		} else {
			now := time.Now().UTC()
			l.VlTimestamp = &now
		}
	}
	return *l.VlTimestamp
}
func (l *logEntryImpl) Level(level string) string {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true, "fatal": true}
	if level != "" && validLevels[level] && l.VlLevel == "" {
		l.VlLevel = LogLevel(level)
	}
	return string(l.VlLevel)
}
func (l *logEntryImpl) Source(source string) string {
	if source != "" && l.VlSource == "" {
		l.VlSource = source
	}
	return l.VlSource
}
func (l *logEntryImpl) Context(context string) string {
	if context != "" && l.VlContext == "" {
		l.VlContext = context
	}
	return l.VlContext
}
func (l *logEntryImpl) Message(message string) string {
	if message != "" && l.VlMessage == "" {
		l.VlMessage = message
	}
	return l.VlMessage
}
func (l *logEntryImpl) Tags(tags map[string]string) map[string]string {
	if l.VlTags == nil {
		l.VlTags = make(map[string]string)
	}
	for k, v := range tags {
		l.VlTags[k] = v
	}
	return l.VlTags
}
func (l *logEntryImpl) Metadata(metadata map[string]any) map[string]any {
	if l.VlMetadata == nil {
		l.VlMetadata = make(map[string]any)
	}
	for k, v := range metadata {
		l.VlMetadata[k] = v
	}
	return l.VlMetadata
}
func (l *logEntryImpl) ProcessID(pid int) int {
	if pid != 0 && l.VlProcessID == 0 {
		l.VlProcessID = pid
	}
	return l.VlProcessID
}
func (l *logEntryImpl) Hostname(hostname string) string {
	if hostname != "" && l.VlHostname == "" {
		l.VlHostname = hostname
	}
	return l.VlHostname
}
func (l *logEntryImpl) Severity(severity int) int {
	if severity > 0 && l.VlSeverity == 0 {
		l.VlSeverity = severity
	}
	return l.VlSeverity
}
func (l *logEntryImpl) TraceID(traceID string) string {
	if traceID != "" && l.VlTraceID == "" {
		l.VlTraceID = traceID
	}
	return l.VlTraceID
}
func (l *logEntryImpl) AddTag(key, value string) {
	if l.VlTags == nil {
		l.VlTags = make(map[string]string)
	}
	l.VlTags[key] = value
}
func (l *logEntryImpl) Validate() error {
	if l.VlTimestamp == nil {
		return fmt.Errorf("timestamp is required")
	}
	if l.VlLevel == "" {
		return fmt.Errorf("level is required")
	}
	if l.VlMessage == "" {
		return fmt.Errorf("message is required")
	}
	if l.VlSeverity <= 0 {
		return fmt.Errorf("severity must be greater than zero")
	}
	return nil
}
func NewLogEntry() LogEntry { return &logEntryImpl{} }
