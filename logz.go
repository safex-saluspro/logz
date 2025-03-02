package main

import (
	"fmt"
	"github.com/faelmori/logz/cmd"
	"os"
)

func main() {
	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		fmt.Printf("Error executing command: %v\n", logzErr)
		os.Exit(1)
	}
}
