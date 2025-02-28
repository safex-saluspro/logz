package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/services"
	"github.com/spf13/cobra"
)

// ServiceCmd retorna um comando cobra para gerenciar o serviço web.
func ServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the web service",
	}
	// Subcomandos: start e stop
	cmd.AddCommand(startServiceCmd())
	cmd.AddCommand(stopServiceCmd())
	return cmd
}

func startServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the web service in the background",
		Run: func(cmd *cobra.Command, args []string) {
			if err := services.Start(); err != nil {
				fmt.Printf("Error starting service: %v\n", err)
			} else {
				fmt.Println("Service started successfully.")
			}
		},
	}
}

func stopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the running web service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := services.Stop(); err != nil {
				fmt.Printf("Error stopping service: %v\n", err)
			} else {
				fmt.Println("Service stopped successfully.")
			}
		},
	}
}
