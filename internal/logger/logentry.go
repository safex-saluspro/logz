package logger

import (
	"errors"
	"fmt"
	"time"
)

// LogLevel representa os níveis de log suportados.
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

// LogEntry concentra os dados de uma entrada de log.
// Observação: cada campo tem um método chainable para facilitar sua construção.
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Source    string    `json:"source"`
	// Embora na versão original _Context_ fosse uma string, o uso de metadados
	// estruturados ficou concentrado em **Metadata**.
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

// NewLogEntry cria uma nova entrada de log com o timestamp atual e inicializa os mapas.
func NewLogEntry() *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Tags:      make(map[string]string),
		Metadata:  make(map[string]interface{}),
	}
}

// WithLevel define o nível do log.
func (le *LogEntry) WithLevel(level LogLevel) *LogEntry {
	le.Level = level
	return le
}

// WithSource define a origem do log.
func (le *LogEntry) WithSource(source string) *LogEntry {
	le.Source = source
	return le
}

// WithContext define um contexto para a entrada (opcional).
func (le *LogEntry) WithContext(context string) *LogEntry {
	le.Context = context
	return le
}

// WithMessage define a mensagem da entrada.
func (le *LogEntry) WithMessage(message string) *LogEntry {
	le.Message = message
	return le
}

// WithProcessID define o ID do processo.
func (le *LogEntry) WithProcessID(pid int) *LogEntry {
	le.ProcessID = pid
	return le
}

// WithHostname define o hostname de onde o log se origina.
func (le *LogEntry) WithHostname(hostname string) *LogEntry {
	le.Hostname = hostname
	return le
}

// WithSeverity define a severidade numérica do log.
func (le *LogEntry) WithSeverity(severity int) *LogEntry {
	le.Severity = severity
	return le
}

// WithTraceID define um id para rastreamento.
func (le *LogEntry) WithTraceID(traceID string) *LogEntry {
	le.TraceID = traceID
	return le
}

// AddTag adiciona ou sobrescreve uma tag.
func (le *LogEntry) AddTag(key, value string) *LogEntry {
	if le.Tags == nil {
		le.Tags = make(map[string]string)
	}
	le.Tags[key] = value
	return le
}

// AddMetadata adiciona ou sobrescreve um metadado.
func (le *LogEntry) AddMetadata(key string, value interface{}) *LogEntry {
	if le.Metadata == nil {
		le.Metadata = make(map[string]interface{})
	}
	le.Metadata[key] = value
	return le
}

// Validate verifica se os campos obrigatórios foram informados.
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

// String retorna uma representação simples da entrada de log.
func (le *LogEntry) String() string {
	return fmt.Sprintf("[%s] %s - %s",
		le.Timestamp.Format(time.RFC3339),
		le.Level,
		le.Message)
}
