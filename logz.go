package main

import (
	"fmt"
	"github.com/faelmori/logz/cmd"
	"github.com/faelmori/logz/internal/services"
	"github.com/spf13/viper"
	"os"
)
import "github.com/faelmori/logz/logger"

func main() {
	// Verifica se o binário foi chamado com o argumento "service-run".
	if len(os.Args) > 2 && os.Args[1] == "service-run" {
		var config *services.Config
		viper.SetConfigFile(os.Args[2])
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Starting service (config file: " + os.Args[2] + ")")
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
		}
		if err := services.Run(); err != nil {
			fmt.Printf("Error running service: %v\n", err)
			os.Exit(1)
		}
	}

	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		fmt.Printf("Error executing command: %v\n", logzErr)
		os.Exit(1)
	}
}

func GetLogger(prefix *string) logger.LoggerInterface { return logger.NewLogger(prefix) }
