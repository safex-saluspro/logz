package logz

import (
	core "github.com/faelmori/logz/internal/logger"
	logz "github.com/faelmori/logz/logger"
	"sync"
)

var (
	prfx   = "Logz"     // Default prefix
	logger LogzLogger   // Global logger instance
	mu     sync.RWMutex // Mutex for concurrency control
	once   sync.Once    // Ensure single initialization
)

type LogLevel = core.LogLevel
type LogFormat = core.LogFormat
type LogzConfig = core.Config
type LogzConfigManager = core.ConfigManager
type LogzNotifierManager = core.NotifierManager
type LogzNotifier = core.Notifier
type LogzLogger = logz.LogzLogger
type LogzWriter = core.LogWriter

// initializeLogger initializes the global logger with the given prefix.
func initializeLogger(prefix string) {
	once.Do(func() {
		if prefix == "" {
			prefix = prfx
		}
		logger = logz.NewLogger(prefix)
	})
}

// GetLogger returns the global logger instance, initializing it if necessary.
func GetLogger(prefix string) LogzLogger {
	initializeLogger(prefix)

	mu.RLock()
	defer mu.RUnlock()
	return logger
}

// NewLogger creates a new logger instance with the given prefix.
func NewLogger(prefix string) LogzLogger {
	if prefix == "" {
		prefix = prfx
	}
	return logz.NewLogger(prefix)
}

// SetLogger sets the global logger instance to the provided logger.
func SetLogger(newLogger LogzLogger) {
	mu.Lock()
	defer mu.Unlock()
	logger = newLogger
}

// SetPrefix sets the global prefix for the logger.
func SetPrefix(prefix string) {
	mu.Lock()
	defer mu.Unlock()
	prfx = prefix
}

// GetPrefix returns the global prefix for the logger.
func GetPrefix() string {
	mu.RLock()
	defer mu.RUnlock()
	return prfx
}

// SetLogLevel sets the log level for the global logger.
func SetLogLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetLevel(level)
	}
}

// GetLogLevel returns the log level of the global logger.
func GetLogLevel() LogLevel {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return core.DEBUG
	}
	return logger.GetLevel()
}

// SetLogWriter sets the log writer for the global logger.
func SetLogWriter(writer LogzWriter) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetWriter(writer)
	}
}

// GetLogWriter returns the log writer of the global logger.
func GetLogWriter() LogzWriter {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.GetWriter()
}

// SetLogConfig sets the configuration for the global logger.
func SetLogConfig(config LogzConfig) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetConfig(config)
	}
}

// GetLogConfig returns the configuration of the global logger.
func GetLogConfig() LogzConfig {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.GetConfig()
}

// SetMetadata sets a metadata key-value pair for the global logger.
func SetMetadata(key string, value interface{}) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.SetMetadata(key, value)
	}
}

// Debug logs a debug message with the given context.
func Debug(msg string, ctx map[string]interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Debug(msg, ctx)
	}
}

// Info logs an info message with the given context.
func Info(msg string, ctx map[string]interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Info(msg, ctx)
	}
}

// Warn logs a warning message with the given context.
func Warn(msg string, ctx map[string]interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Warn(msg, ctx)
	}
}

// Error logs an error message with the given context.
func Error(msg string, ctx map[string]interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.Error(msg, ctx)
	}
}

// FatalC logs a fatal message with the given context and exits the application.
func FatalC(msg string, ctx map[string]interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if logger != nil {
		logger.FatalC(msg, ctx)
	}
}

// AddNotifier adds a notifier to the global logger's configuration.
func AddNotifier(name string, notifier LogzNotifier) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.
			GetConfig().
			NotifierManager().
			AddNotifier(name, notifier)
	}
}

// GetNotifier returns the notifier with the given name from the global logger's configuration.
func GetNotifier(name string) (LogzNotifier, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return nil, false
	}
	return logger.
		GetConfig().
		NotifierManager().
		GetNotifier(name)
}

// ListNotifiers returns a list of all notifier names in the global logger's configuration.
func ListNotifiers() []string {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return nil
	}
	return logger.
		GetConfig().
		NotifierManager().
		ListNotifiers()
}

// SetLogFormat sets the log format for the global logger.
func SetLogFormat(format LogFormat) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.
			GetConfig().
			SetFormat(format)
	}
}

// GetLogFormat returns the log format of the global logger.
func GetLogFormat() string {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return "text"
	}
	return logger.GetConfig().Format()
}

// SetLogOutput sets the log output for the global logger.
func SetLogOutput(output string) {
	mu.Lock()
	defer mu.Unlock()
	if logger != nil {
		logger.GetConfig().SetOutput(output)
	}
}

// GetLogOutput returns the log output of the global logger.
func GetLogOutput() string {
	mu.RLock()
	defer mu.RUnlock()
	if logger == nil {
		return "stdout"
	}
	return logger.GetConfig().Output()
}
