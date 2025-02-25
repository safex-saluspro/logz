package utils

import "time"

var (
	MaxLogSize int64 = 60 * 1024 * 1024 // 60 MB
	LogModule  string
	LogLevel   string
	LogOutput  string
	LogColor   = map[string]string{
		"debug":   "\033[34m",
		"info":    "\033[36m",
		"warn":    "\033[33m",
		"error":   "\033[31m",
		"success": "\033[32m",
		"answer":  "\033[35m",
		"default": "\033[0m",
	}
	LogLevels = map[string]int{
		"DEBUG": 1,
		"INFO":  2,
		"WARN":  3,
		"ERROR": 4,
	}
)

type Options struct {
	logType     string
	message     string
	whatToShow  string
	follow      bool
	whatToClear string
}

type LogMetrics struct {
	InfoCount    int `json:"info_count"`
	WarnCount    int `json:"warn_count"`
	ErrorCount   int `json:"error_count"`
	DebugCount   int `json:"debug_count"`
	SuccessCount int `json:"success_count"`
}

type LogEntryImpl struct {
	Timestamp time.Time         `json:"timestamp"`          // Timestamp padronizado UTC
	Level     string            `json:"level"`              // Nível do log: debug, info, warn, error, fatal
	Source    string            `json:"source"`             // Origem: sistema, aplicação, processo específico
	Context   string            `json:"context"`            // Contexto do log (ex: app específica, microserviço, módulo)
	Message   string            `json:"message"`            // Mensagem principal do log
	Tags      map[string]string `json:"tags"`               // Tags para categorização (ex: "host": "server01")
	Metadata  map[string]any    `json:"metadata"`           // Dados adicionais opcionais
	ProcessID int               `json:"pid,omitempty"`      // ID do processo (se aplicável)
	Hostname  string            `json:"hostname"`           // Nome da máquina (para logs de SO e distribuídos)
	Severity  int               `json:"severity"`           // Severidade numérica (compatível com syslog)
	TraceID   string            `json:"trace_id,omitempty"` // Para tracing distribuído
}
