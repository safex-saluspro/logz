package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogMode string
type LogFormat string

const (
	JSONFormat     LogFormat = "json"
	TextFormat     LogFormat = "text"
	ModeService    LogMode   = "service"    // Indica que o logger está sendo usado por um processo destacado
	ModeStandalone LogMode   = "standalone" // Indica que o logger está sendo usado localmente (ex.: CLI)
)

var logLevels = map[LogLevel]int{
	DEBUG: 1,
	INFO:  2,
	WARN:  3,
	ERROR: 4,
	FATAL: 5,
}

type Logger struct {
	level    LogLevel
	writer   LogWriter
	config   Config
	metadata map[string]interface{}
	mode     LogMode // Controle de modo: service ou standalone
}

func NewLogger(config Config) *Logger {
	// Define o nível de log a partir do Config
	level := LogLevel(config.Level()) // Método config.Level() retorna o nível de log como string

	var out *os.File
	if config.DefaultLogPath() == "stdout" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile(config.DefaultLogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Erro ao abrir arquivo de log: %v\nRedirecionando para stdout...\n", err)
			out = os.Stdout
		}
	}

	// Inicializa o formatador (JSON ou texto)
	var formatter LogFormatter
	if config.Format() == "json" {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{}
	}
	writer := NewDefaultWriter(out, formatter)

	// Lê o modo do Config
	mode := config.Mode()
	if mode != ModeService && mode != ModeStandalone {
		mode = ModeStandalone // Default para standalone se não especificado
	}

	return &Logger{
		level:    level,
		writer:   writer,
		config:   config,
		metadata: make(map[string]interface{}),
		mode:     mode,
	}
}

func (l *Logger) SetMetadata(key string, value interface{}) {
	l.metadata[key] = value
}
func (l *Logger) shouldLog(level LogLevel) bool {
	return logLevels[level] >= logLevels[l.level]
}
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

	// Mescla os metadados globais e locais
	finalContext := mergeContext(l.metadata, ctx)
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	// Escreve o log utilizando o writer configurado
	if err := l.writer.Write(entry); err != nil {
		log.Printf("Erro ao escrever log: %v", err)
	}

	// Apenas no modo de serviço, notificar via Notifiers
	if l.mode == ModeService && l.config != nil {
		for _, name := range l.config.NotifierManager().ListNotifiers() {
			if notifier, ok := l.config.NotifierManager().GetNotifier(name); ok {
				if notifier != nil {
					ntf := *notifier
					if ntfErr := ntf.Notify(entry); ntfErr != nil {
						log.Printf("Erro ao notificar %s: %v", name, ntfErr)
					}
				}
			}
		}
	}

	// Atualiza métricas no PrometheusManager, se habilitado
	if l.mode == ModeService {
		pm := GetPrometheusManager()
		if pm.IsEnabled() {
			pm.IncrementMetric("logs_total", 1)
			pm.IncrementMetric("logs_total_"+string(level), 1)
		}
	}

	// Finaliza o processo em caso de log FATAL
	if level == FATAL {
		os.Exit(1)
	}
}
func (l *Logger) Debug(msg string, ctx map[string]interface{}) { l.log(DEBUG, msg, ctx) }
func (l *Logger) Info(msg string, ctx map[string]interface{})  { l.log(INFO, msg, ctx) }
func (l *Logger) Warn(msg string, ctx map[string]interface{})  { l.log(WARN, msg, ctx) }
func (l *Logger) Error(msg string, ctx map[string]interface{}) { l.log(ERROR, msg, ctx) }
func (l *Logger) Fatal(msg string, ctx map[string]interface{}) { l.log(FATAL, msg, ctx) }

func getCallerInfo(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", trimFilePath(file), line, funcName)
}
func trimFilePath(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return filePath
}
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
