package logger

import (
	"github.com/faelmori/logz/internal/logger"
	"log"
	"os"
)

// LogzCore é a interface com os métodos básicos do logger existente.
type LogzCore interface {
	SetMetadata(key string, value interface{})
	Debug(msg string, ctx map[string]interface{})
	Info(msg string, ctx map[string]interface{})
	Warn(msg string, ctx map[string]interface{})
	Error(msg string, ctx map[string]interface{})
	FatalC(msg string, ctx map[string]interface{})
}

// LogzLogger combina o logger existente com os métodos padrões de log do Go.
type LogzLogger interface {
	LogzCore
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

// logzLogger é a implementação do LoggerInterface, unificando o novo LogzCoreImpl e o antigo.
type logzLogger struct {
	logger     *log.Logger
	coreLogger *logger.LogzCoreImpl
}

// Métodos do logger Go padrão (compatível com log.Logger)
func (l *logzLogger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

func (l *logzLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *logzLogger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *logzLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}

func (l *logzLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func (l *logzLogger) Fatalln(v ...interface{}) {
	l.logger.Fatalln(v...)
}

func (l *logzLogger) Panic(v ...interface{}) {
	l.logger.Panic(v...)
}

func (l *logzLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

func (l *logzLogger) Panicln(v ...interface{}) {
	l.logger.Panicln(v...)
}

func (l *logzLogger) Debug(msg string, ctx map[string]interface{}) {
	l.coreLogger.Debug(msg, ctx)
}

func (l *logzLogger) Info(msg string, ctx map[string]interface{}) {
	l.coreLogger.Info(msg, ctx)
}

func (l *logzLogger) Warn(msg string, ctx map[string]interface{}) {
	l.coreLogger.Warn(msg, ctx)
}

func (l *logzLogger) Error(msg string, ctx map[string]interface{}) {
	l.coreLogger.Error(msg, ctx)
}

func (l *logzLogger) FatalC(msg string, ctx map[string]interface{}) {
	l.coreLogger.FatalC(msg, ctx)
}

func (l *logzLogger) SetMetadata(key string, value interface{}) {
	l.coreLogger.SetMetadata(key, value)
}

func NewLogger(prefix *string) LogzLogger {
	return &logzLogger{
		logger: log.New(
			os.Stdout,
			func() string {
				if prefix != nil {
					return *prefix
				}
				return ""
			}(),
			log.LstdFlags,
		),
	}
}
