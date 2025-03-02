package main

import (
	"fmt"
	"github.com/faelmori/logz/cmd"
	"github.com/faelmori/logz/logger"
	"os"
)

func main() {
	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		fmt.Printf("Error executing command: %v\n", logzErr)
		os.Exit(1)
	}
}

func New(prefix string) logger.LogzLogger {
	return logger.NewLogger(prefix)
}
