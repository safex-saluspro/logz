package logger

import (
	"log"
	"os"
)

// Mapping dos níveis de log para um valor numérico de severidade.
var logLevels = map[LogLevel]int{
	DEBUG: 1,
	INFO:  2,
	WARN:  3,
	ERROR: 4,
	FATAL: 5,
}

// Logger orquestra a criação da entrada de log, sua escrita e o envio para notifiers.
type Logger struct {
	level     LogLevel
	writer    LogWriter
	notifiers []Notifier
	metadata  map[string]interface{}
}

// NewLogger cria uma nova instância de Logger com base nos parâmetros fornecidos.
func NewLogger(level LogLevel, format string, outputPath, externalURL, zmqEndpoint, discordWebhook string) *Logger {
	var out *os.File
	if outputPath == "stdout" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Erro ao abrir arquivo de log: %v", err)
		}
	}

	var formatter LogFormatter
	if format == "json" {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{}
	}
	writer := NewDefaultWriter(out, formatter)

	notifiers := []Notifier{}
	if externalURL != "" || zmqEndpoint != "" {
		extNotifier := NewExternalNotifier(externalURL, zmqEndpoint)
		notifiers = append(notifiers, extNotifier)
	}
	if discordWebhook != "" {
		discNotifier := NewDiscordNotifier(discordWebhook)
		notifiers = append(notifiers, discNotifier)
	}

	return &Logger{
		level:     level,
		writer:    writer,
		notifiers: notifiers,
		metadata:  make(map[string]interface{}),
	}
}

// ParseLogLevel converte uma string para LogLevel; valores inválidos retornam INFO.
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// SetMetadata adiciona um metadado global que será mesclado com o contexto de cada log.
func (l *Logger) SetMetadata(key string, value interface{}) {
	l.metadata[key] = value
}

func (l *Logger) shouldLog(level LogLevel) bool {
	return logLevels[level] >= logLevels[l.level]
}

// log cria uma entrada de log, mescla os metadados e delega a escrita e envio.
func (l *Logger) log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	// Mescla metadados globais com o contexto informado.
	finalContext := mergeContext(l.metadata, ctx)
	// Construção da entrada de log utilizando os métodos chainable.
	entry := NewLogEntry().
		WithLevel(level).
		WithMessage(msg).
		WithSeverity(logLevels[level])
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	if err := l.writer.Write(entry); err != nil {
		log.Printf("Erro ao escrever log: %v", err)
	}
	for _, notifier := range l.notifiers {
		notifier.Notify(entry)
	}
	if level == FATAL {
		os.Exit(1)
	}
}

func (l *Logger) Debug(msg string, ctx map[string]interface{}) { l.log(DEBUG, msg, ctx) }
func (l *Logger) Info(msg string, ctx map[string]interface{})  { l.log(INFO, msg, ctx) }
func (l *Logger) Warn(msg string, ctx map[string]interface{})  { l.log(WARN, msg, ctx) }
func (l *Logger) Error(msg string, ctx map[string]interface{}) { l.log(ERROR, msg, ctx) }
func (l *Logger) Fatal(msg string, ctx map[string]interface{}) { l.log(FATAL, msg, ctx) }

// mergeContext une os metadados globais e o contexto específico.
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
