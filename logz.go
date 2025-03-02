package main

import (
	"fmt"
	"github.com/faelmori/logz/cmd"
	"os"
)
import "github.com/faelmori/logz/logger"

func main() {
	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		fmt.Printf("Error executing command: %v\n", logzErr)
		os.Exit(1)
	}
}

func GetLogger(prefix *string) logger.LogzLogger { return logger.NewLogger(prefix) }
