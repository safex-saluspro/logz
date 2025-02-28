package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/pebbe/zmq4"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type LogFormat string

const (
	JSONFormat LogFormat = "json"
	TextFormat LogFormat = "text"
)

type Logger struct {
	level          LogLevel
	output         *os.File
	metadata       map[string]interface{}
	format         LogFormat
	externalURL    string
	zmqSocket      *zmq4.Socket
	discordWebhook string
}

// NewLogger cria uma nova instância do Logger
func NewLogger(level LogLevel, format LogFormat, outputPath, externalURL, zmqEndpoint, discordWebhook string) *Logger {
	var output *os.File
	if outputPath == "stdout" {
		output = os.Stdout
	} else {
		var err error
		output, err = os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Erro ao abrir arquivo de log: %v", err)
		}
	}

	var zmqSocket *zmq4.Socket
	if zmqEndpoint != "" {
		zmqSocket, _ = zmq4.NewSocket(zmq4.PUSH)
		zmqSocket.Connect(zmqEndpoint)
	}

	return &Logger{
		level:          level,
		output:         output,
		metadata:       make(map[string]interface{}),
		format:         format,
		externalURL:    externalURL,
		zmqSocket:      zmqSocket,
		discordWebhook: discordWebhook,
	}
}

// SetMetadata permite adicionar metadados globais no logger
func (l *Logger) SetMetadata(key string, value interface{}) {
	l.metadata[key] = value
}

// Log escreve uma entrada de log formatada
func (l *Logger) Log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}

	entry := NewLogEntry()

	l.writeLog(entry)
	l.sendToExternal(entry)
	l.sendToDiscord(entry)
}

// Métodos de conveniência
func (l *Logger) Debug(msg string, ctx map[string]interface{}) { l.Log(DEBUG, msg, ctx) }
func (l *Logger) Info(msg string, ctx map[string]interface{})  { l.Log(INFO, msg, ctx) }
func (l *Logger) Warn(msg string, ctx map[string]interface{})  { l.Log(WARN, msg, ctx) }
func (l *Logger) Error(msg string, ctx map[string]interface{}) { l.Log(ERROR, msg, ctx) }
func (l *Logger) Fatal(msg string, ctx map[string]interface{}) {
	l.Log(FATAL, msg, ctx)
	os.Exit(1)
}

func (l *Logger) shouldLog(level LogLevel) bool {
	logLevels := map[LogLevel]int{DEBUG: 1, INFO: 2, WARN: 3, ERROR: 4, FATAL: 5}
	return logLevels[level] >= logLevels[l.level]
}

func (l *Logger) writeLog(entry LogEntry) {
	var output string
	if l.format == JSONFormat {
		data, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Erro ao serializar log: %v", err)
			return
		}
		output = string(data)
	} else {
		output = fmt.Sprintf("[%s] %s - %s (%s)\n", entry.Timestamp, entry.Level, entry.Message, entry.Context)
	}
	fmt.Fprintln(l.output, output)
}

func (l *Logger) sendToExternal(entry LogEntry) {
	if l.externalURL != "" {
		data, _ := json.Marshal(entry)
		_, err := http.Post(l.externalURL, "application/json", strings.NewReader(string(data)))
		if err != nil {
			log.Printf("Erro ao enviar log para %s: %v", l.externalURL, err)
		}
	}

	if l.zmqSocket != nil {
		data, _ := json.Marshal(entry)
		l.zmqSocket.Send(string(data), 0)
	}
}

func (l *Logger) sendToDiscord(entry LogEntry) {
	if l.discordWebhook == "" {
		return
	}
	message := fmt.Sprintf("**[%s] %s**\n%s", entry.Level, entry.Timestamp, entry.Message)
	payload := map[string]string{"content": message}
	jsonPayload, _ := json.Marshal(payload)
	_, err := http.Post(l.discordWebhook, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		log.Printf("Erro ao enviar log para Discord: %v", err)
	}
}

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
