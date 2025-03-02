package logz

import (
	"github.com/faelmori/logz/logger"
)

func New(prefix string) logger.LogzLogger {
	return logger.NewLogger(prefix)
}
