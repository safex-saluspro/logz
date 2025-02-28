package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
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

// Format converte a entrada de log para uma string formatada com cores e ícones.
func (f *TextFormatter) Format(entry *LogEntry) (string, error) {
	// Se a variável LOGZ_NO_COLOR está definida ou se estamos em Windows, desliga cores.
	noColor := os.Getenv("LOGZ_NO_COLOR") != "" || runtime.GOOS == "windows"

	var icon, levelStr, reset string
	if noColor {
		icon = ""
		levelStr = string(entry.Level)
		reset = ""
	} else {
		// Define o "reset" ANSI
		reset = "\033[0m"
		var color string
		// Escolhe cor e ícone conforme o nível
		switch entry.Level {
		case DEBUG:
			color = "\033[34m" // azul
			icon = "🐛"
		case INFO:
			color = "\033[32m" // verde
			icon = "ℹ️"
		case WARN:
			color = "\033[33m" // amarelo
			icon = "⚠️"
		case ERROR:
			color = "\033[31m" // vermelho
			icon = "❌"
		case FATAL:
			color = "\033[35m" // magenta
			icon = "💀"
		default:
			color = ""
			icon = ""
		}
		icon = color + icon + reset
		levelStr = color + string(entry.Level) + reset
	}

	// A formatação inclui timestamp, ícone, nível, mensagem e contexto.
	return fmt.Sprintf("[%s] %s %s - %s (%s)",
		entry.Timestamp.Format(time.RFC3339),
		icon,
		levelStr,
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
