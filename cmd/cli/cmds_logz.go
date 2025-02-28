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

// LogzCmds retorna os comandos CLI para os diferentes níveis de log.
func LogzCmds() []*cobra.Command {
	return []*cobra.Command{
		newLogCmd("debug", []string{"dbg"}),
		newLogCmd("info", []string{"inf"}),
		newLogCmd("warn", []string{"wrn"}),
		newLogCmd("error", []string{"err"}),
		newLogCmd("fatal", []string{"ftl"}),
		watchLogsCmd(),
	}
}

// newLogCmd configura um comando para um nível de log específico.
func newLogCmd(level string, aliases []string) *cobra.Command {
	var format, outputPath, externalURL, zmqEndpoint, discordWebhook string
	var metaData, ctx map[string]string
	var msg string

	cmd := &cobra.Command{
		Use:     level,
		Aliases: aliases,
		Short:   "Loga uma mensagem de nível " + level,
		Run: func(cmd *cobra.Command, args []string) {
			// Cria o Logger usando os parâmetros fornecidos.
			logr := logger.NewLogger(logger.ParseLogLevel(level), format, outputPath, externalURL, zmqEndpoint, discordWebhook)
			for k, v := range metaData {
				logr.SetMetadata(k, v)
			}
			// Converte o contexto de string para map[string]interface{}
			ctxInterface := make(map[string]interface{})
			for k, v := range ctx {
				ctxInterface[k] = v
			}
			// Chama o método de log conforme o nível informado.
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

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Formato do log (text ou json)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "stdout", "Destino do output do log")
	cmd.Flags().StringVarP(&externalURL, "external-url", "e", "", "URL externa para enviar o log")
	cmd.Flags().StringVarP(&zmqEndpoint, "zmq-endpoint", "z", "", "Endpoint ZMQ")
	cmd.Flags().StringVarP(&discordWebhook, "discord-webhook", "d", "", "Webhook do Discord")
	cmd.Flags().StringToStringVarP(&metaData, "metadata", "m", nil, "Metadados a serem incluídos")
	cmd.Flags().StringToStringVarP(&ctx, "context", "c", nil, "Contexto para o log")
	cmd.Flags().StringVarP(&msg, "msg", "M", "", "Mensagem do log")

	return cmd
}

func watchLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Watch logs in real time",
		Run: func(cmd *cobra.Command, args []string) {
			logFilePath := "./logz.log" // Ajuste isso conforme sua configuração
			reader := logger.NewFileLogReader()
			stopChan := make(chan struct{})
			// Captura sinais para interrupção
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				close(stopChan)
			}()
			fmt.Println("Watching logs (press Ctrl+C to exit):")
			if err := reader.Tail(logFilePath, stopChan); err != nil {
				fmt.Printf("Error watching logs: %v\n", err)
			}
			// Aguarda um pequeno delay para finalizar
			time.Sleep(500 * time.Millisecond)
		},
	}
}
