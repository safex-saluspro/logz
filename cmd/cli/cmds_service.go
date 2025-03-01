package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/logger"
	"github.com/spf13/cobra"
)

func ServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the web service",
	}
	cmd.AddCommand(startServiceCmd())
	cmd.AddCommand(stopServiceCmd())
	cmd.AddCommand(getServiceCmd())
	cmd.AddCommand(spawnServiceCmd())
	return cmd
}

func startServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Inicia o serviço destacado",
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

			if err := logger.Start(config.Port()); err != nil {
				fmt.Printf("Erro ao iniciar serviço: %v\n", err)
			} else {
				fmt.Println("Serviço iniciado com sucesso.")
			}
		},
	}
}
func stopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Para o serviço destacado",
		Run: func(cmd *cobra.Command, args []string) {
			if err := logger.Stop(); err != nil {
				fmt.Printf("Erro ao parar serviço: %v\n", err)
			} else {
				fmt.Println("Serviço parado com sucesso.")
			}
		},
	}
}
func getServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get information about the running service",
		Run: func(cmd *cobra.Command, args []string) {
			pid, port, pidPath, err := logger.GetServiceInfo()
			if err != nil {
				fmt.Println("Service is not running")
			} else {
				fmt.Printf("Service running with PID %d on port %s\n", pid, port)
				fmt.Printf("PID file: %s\n", pidPath)
			}
		},
	}
}
func spawnServiceCmd() *cobra.Command {
	var configPath string
	spCmd := &cobra.Command{
		Use:    "spawn",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := logger.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	spCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the service configuration file")
	return spCmd
}
