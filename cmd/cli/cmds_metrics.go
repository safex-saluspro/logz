package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/services"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func MetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Manage Prometheus metrics",
	}

	cmd.AddCommand(enableMetricsCmd())
	cmd.AddCommand(disableMetricsCmd())
	cmd.AddCommand(addMetricCmd())
	cmd.AddCommand(removeMetricCmd())
	cmd.AddCommand(listMetricsCmd())
	cmd.AddCommand(watchMetricsCmd())

	return cmd
}

func enableMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "enable",
		Aliases: []string{"en"},
		Short:   "Enable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			pm.Enable()
		},
	}
}

func disableMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "disable",
		Aliases: []string{"dis"},
		Short:   "Disable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			pm.Disable()
		},
	}
}

func addMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "add [name] [value]",
		Aliases: []string{"a"},
		Short:   "Add or update a Prometheus metric",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			value, valueErr := strconv.ParseFloat(args[1], 64)
			if valueErr != nil {
				fmt.Printf("Invalid metric value: %v\n", valueErr)
				return
			}
			pm := services.GetPrometheusManager()
			pm.AddMetric(name, value)
		},
	}
}

func removeMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove [name]",
		Aliases: []string{"r"},
		Short:   "Remove a Prometheus metric",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			pm := services.GetPrometheusManager()
			pm.RemoveMetric(name)
		},
	}
}

func listMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "List all Prometheus metrics",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			metrics := pm.GetMetrics()
			if len(metrics) == 0 {
				fmt.Println("No metrics registered.")
				return
			}
			fmt.Println("Registered metrics:")
			for name, value := range metrics {
				fmt.Printf(" - %s: %f\n", name, value)
			}
		},
	}
}

func watchMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Short:   "Watch Prometheus metrics in real time",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Watching metrics (press Ctrl+C to exit):")
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					metrics := services.GetPrometheusManager().GetMetrics()
					fmt.Println("Current Metrics:")
					if len(metrics) == 0 {
						fmt.Println("  No metrics registered.")
					} else {
						for name, value := range metrics {
							fmt.Printf(" - %s: %f\n", name, value)
						}
					}
					fmt.Println("-----")
				}
			}
		},
	}
}
