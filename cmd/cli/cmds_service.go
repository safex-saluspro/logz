package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/services"
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
	var port string

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the web service in the background",
		Run: func(cmd *cobra.Command, args []string) {
			if err := services.Start(port); err != nil {
				fmt.Printf("Error starting service: %v\n", err)
			}
		},
	}

	startCmd.Flags().StringVarP(&port, "port", "p", "9999", "Port to listen on")
	return startCmd
}
func stopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the running web service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := services.Stop(); err != nil {
				fmt.Printf("Error stopping service: %v\n", err)
			}
		},
	}
}
func getServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Get information about the running service",
		Run: func(cmd *cobra.Command, args []string) {
			pid, port, pidPath, err := services.GetServiceInfo()
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
			if err := services.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	spCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the service configuration file")
	return spCmd
}
