package logger

import (
	"log"
	"os"
)

// ExistingLoggerInterface is the interface you already use
type ExistingLoggerInterface interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
}

// LoggerInterface is the new interface that combines both
type LoggerInterface interface {
	ExistingLoggerInterface
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type logzLogger struct {
	logger *log.Logger
}

func (l *logzLogger) Info(args ...interface{}) {
	l.logger.SetPrefix("INFO: ")
	l.logger.Println(args...)
}
func (l *logzLogger) Warn(args ...interface{}) {
	l.logger.SetPrefix("WARN: ")
	l.logger.Println(args...)
}
func (l *logzLogger) Error(args ...interface{}) {
	l.logger.SetPrefix("ERROR: ")
	l.logger.Println(args...)
}
func (l *logzLogger) Debug(args ...interface{}) {
	l.logger.SetPrefix("DEBUG: ")
	l.logger.Println(args...)
}
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

func NewLogger(prefix *string) LoggerInterface {
	return &logzLogger{
		logger: log.New(
			os.Stdout,
			func() string {
				if prefix != nil {
					return *prefix
				} else {
					return ""
				}
			}(),
			log.LstdFlags),
	}
}
