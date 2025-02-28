package main

import (
	"github.com/faelmori/logz/cmd"
	"github.com/faelmori/logz/internal/services"
	"os"
)
import "github.com/faelmori/logz/logger"

func main() {
	// Verifica se o binário foi chamado com o argumento "service-run".
	if len(os.Args) > 1 && os.Args[1] == "service-run" {
		// Executa o serviço (este processo bloqueia enquanto o serviço estiver ativo)
		if err := services.Run(); err != nil {
			os.Exit(1)
		}
		return
	}

	if logzErr := cmd.RegX().Execute(); logzErr != nil {
		panic(logzErr)
	}
}

func GetLogger(prefix *string) logger.LoggerInterface { return logger.NewLogger(prefix) }
