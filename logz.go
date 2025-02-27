package main

import "github.com/faelmori/logz/cmd"
import "github.com/faelmori/logz/logger"

func main() {
	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		panic(logzErr)
	}
}

func GetLogger(prefix *string) logger.LoggerInterface { return logger.NewLogger(prefix) }
