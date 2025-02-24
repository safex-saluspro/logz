package cmd

var (
	maxLogSize int64 = 60 * 1024 * 1024 // 60 MB
	LogModule  string
	logLevel   string
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
