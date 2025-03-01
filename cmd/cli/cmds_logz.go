package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/logger"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// LogzCmds retorna os comandos CLI para diferentes níveis de log e gerenciamento.
func LogzCmds() []*cobra.Command {
	return []*cobra.Command{
		newLogCmd("debug", []string{"dbg"}),
		newLogCmd("info", []string{"inf"}),
		newLogCmd("warn", []string{"wrn"}),
		newLogCmd("error", []string{"err"}),
		newLogCmd("fatal", []string{"ftl"}),
		watchLogsCmd(),
		startServiceCmd(),
		stopServiceCmd(),
	}
}

// newLogCmd configura um comando para um nível de log específico.
func newLogCmd(level string, aliases []string) *cobra.Command {
	var metaData, ctx map[string]string
	var msg string

	cmd := &cobra.Command{
		Use:     level,
		Aliases: aliases,
		Short:   "Loga uma mensagem de nível " + level,
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Erro ao inicializar ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Erro ao carregar configuração: %v\n", err)
				return
			}

			logr := logger.NewLogger(config)
			for k, v := range metaData {
				logr.SetMetadata(k, v)
			}
			ctxInterface := make(map[string]interface{})
			for k, v := range ctx {
				ctxInterface[k] = v
			}
			switch level {
			case "debug":
				logr.Debug(msg, ctxInterface)
			case "info":
				logr.Info(msg, ctxInterface)
			case "warn":
				logr.Warn(msg, ctxInterface)
			case "error":
				logr.Error(msg, ctxInterface)
			case "fatal":
				logr.Fatal(msg, ctxInterface)
			default:
				logr.Info(msg, ctxInterface)
			}
		},
	}

	cmd.Flags().StringToStringVarP(&metaData, "metadata", "m", nil, "Metadados a serem incluídos")
	cmd.Flags().StringToStringVarP(&ctx, "context", "c", nil, "Contexto para o log")
	cmd.Flags().StringVarP(&msg, "msg", "M", "", "Mensagem do log")

	return cmd
}

// watchLogsCmd monitora logs em tempo real.
func watchLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Monitora logs em tempo real",
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Erro ao inicializar ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Erro ao carregar configuração: %v\n", err)
				return
			}

			logFilePath := config.DefaultLogPath()
			reader := logger.NewFileLogReader()
			stopChan := make(chan struct{})

			// Captura sinais para interrupção
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				close(stopChan)
			}()

			fmt.Println("Monitoração iniciada (Ctrl+C para sair):")
			if err := reader.Tail(logFilePath, stopChan); err != nil {
				fmt.Printf("Erro ao monitorar logs: %v\n", err)
			}

			// Aguarda um pequeno delay para finalizar
			time.Sleep(500 * time.Millisecond)
		},
	}
}

// parseLogLevel mapeia o nível de log para o enum correspondente.
func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.DEBUG
	case "info":
		return logger.INFO
	case "warn":
		return logger.WARN
	case "error":
		return logger.ERROR
	case "fatal":
		return logger.FATAL
	default:
		return logger.INFO
	}
}
