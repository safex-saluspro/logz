package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogFormatter define o contrato para formatação de entradas de log.
type LogFormatter interface {
	Format(entry *LogEntry) (string, error)
}

// JSONFormatter formata o log no padrão JSON.
type JSONFormatter struct{}

// Format converte a entrada de log para JSON.
func (f *JSONFormatter) Format(entry *LogEntry) (string, error) {
	data, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// TextFormatter formata o log em texto simples.
type TextFormatter struct{}

// Format converte a entrada de log para uma string formata.
func (f *TextFormatter) Format(entry *LogEntry) (string, error) {
	return fmt.Sprintf("[%s] %s - %s (%s)",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		entry.Context,
	), nil
}

// LogWriter define o contrato para escrita de log.
type LogWriter interface {
	Write(entry *LogEntry) error
}

// DefaultWriter implementa LogWriter utilizando um io.Writer e um LogFormatter.
type DefaultWriter struct {
	out       io.Writer
	formatter LogFormatter
}

// NewDefaultWriter cria uma nova instância de DefaultWriter.
func NewDefaultWriter(out io.Writer, formatter LogFormatter) *DefaultWriter {
	return &DefaultWriter{
		out:       out,
		formatter: formatter,
	}
}

// Write formata a entrada e a escreve no destino configurado.
func (w *DefaultWriter) Write(entry *LogEntry) error {
	formatted, err := w.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w.out, formatted)
	return err
}
