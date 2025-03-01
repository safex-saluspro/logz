package logger

import (
	"fmt"
	"github.com/faelmori/logz/internal/services"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
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
}

func NewLogger(level LogLevel, format, outputPath string) *Logger {
	var out *os.File
	if outputPath == "stdout" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(fmt.Sprintf("Erro ao abrir arquivo de log: %v\nRedirecionando para stdout...", err))
			out = os.Stdout
		}
	}

	var formatter LogFormatter
	if format == "json" {
		formatter = &JSONFormatter{}
	} else {
		formatter = &TextFormatter{}
	}
	writer := NewDefaultWriter(out, formatter)
	configManager := NewConfigManager()

	var config Config
	var configErr error

	if configManager != nil {
		cfgMgr := *configManager
		config, configErr = cfgMgr.LoadConfig()
		if configErr != nil {
			log.Fatalf("Erro ao carregar configuração: %v", configErr)
		}
	}

	return &Logger{
		level:    level,
		writer:   writer,
		config:   config,
		metadata: make(map[string]interface{}),
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
	finalContext := mergeContext(l.metadata, ctx)
	for k, v := range finalContext {
		entry.AddMetadata(k, v)
	}

	if err := l.writer.Write(entry); err != nil {
		log.Printf("Erro ao escrever log: %v", err)
	}

	if l.config != nil {
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

	pm := services.GetPrometheusManager()
	if pm.IsEnabled() {
		pm.IncrementMetric("logs_total", 1)
		pm.IncrementMetric("logs_total_"+string(level), 1)
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
