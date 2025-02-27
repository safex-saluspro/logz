package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type LoggerLogz struct {
	level    LogLevel
	output   *os.File
	metadata map[string]interface{}
}

// SetMetadata permite adicionar metadados globais no logger
func (l *LoggerLogz) SetMetadata(key string, value interface{}) {
	l.metadata[key] = value
}

// Log escreve uma entrada de log formatada
func (l *LoggerLogz) Log(level LogLevel, msg string, ctx map[string]interface{}) {
	if !l.shouldLog(level) {
		return
	}
	timestamp := time.Now().UTC()
	caller := getCallerInfo(3)

	entry := LogRegistry{
		Timestamp: timestamp,
		Level:     level,
		Message:   msg,
		Caller:    caller,
		Context:   mergeContext(l.metadata, ctx),
	}
	l.writeLog(entry)
}

// Métodos de conveniência
func (l *LoggerLogz) Debug(msg string, ctx map[string]interface{}) { l.Log(DEBUG, msg, ctx) }
func (l *LoggerLogz) Info(msg string, ctx map[string]interface{})  { l.Log(INFO, msg, ctx) }
func (l *LoggerLogz) Warn(msg string, ctx map[string]interface{})  { l.Log(WARN, msg, ctx) }
func (l *LoggerLogz) Error(msg string, ctx map[string]interface{}) { l.Log(ERROR, msg, ctx) }
func (l *LoggerLogz) Fatal(msg string, ctx map[string]interface{}) {
	l.Log(FATAL, msg, ctx)
	os.Exit(1)
}

// shouldLog verifica se o log deve ser registrado
func (l *LoggerLogz) shouldLog(level LogLevel) bool {
	logLevels := map[LogLevel]int{DEBUG: 1, INFO: 2, WARN: 3, ERROR: 4, FATAL: 5}
	return logLevels[level] >= logLevels[l.level]
}
func (l *LoggerLogz) writeLog(entry LogRegistry) {
	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Erro ao serializar log: %v", err)
		return
	}
	_, err = fmt.Fprintln(l.output, string(data))
	if err != nil {
		return
	}
}

// NewLogger cria uma nova instância do LoggerLogz
func OldLogger(level LogLevel, outputPath string) *LoggerLogz {
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
	return &LoggerLogz{
		level:    level,
		output:   output,
		metadata: make(map[string]interface{}),
	}
}

//func getCallerInfo(skip int) string {
//	pc, file, line, ok := runtime.Caller(skip)
//	if !ok {
//		return "unknown"
//	}
//	funcName := runtime.FuncForPC(pc).Name()
//	return fmt.Sprintf("%s:%d %s", trimFilePath(file), line, funcName)
//}
//func trimFilePath(filePath string) string {
//	parts := strings.Split(filePath, "/")
//	if len(parts) > 2 {
//		return strings.Join(parts[len(parts)-2:], "/")
//	}
//	return filePath
//}
//func mergeContext(global, local map[string]interface{}) map[string]interface{} {
//	merged := make(map[string]interface{})
//	for k, v := range global {
//		merged[k] = v
//	}
//	for k, v := range local {
//		merged[k] = v
//	}
//	return merged
//}
