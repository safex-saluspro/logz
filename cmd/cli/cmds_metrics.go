package cli

import (
	"github.com/faelmori/logz/internal/services"
	"github.com/spf13/cobra"
	"strconv"
)

// MetricsCmd retorna os comandos relacionados ao Prometheus e métricas.
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

	return cmd
}

func enableMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			pm.Enable()
		},
	}
}

func disableMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable Prometheus integration",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			pm.Disable()
		},
	}
}

func addMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [name] [value]",
		Short: "Add or update a Prometheus metric",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			value, valueErr := strconv.ParseFloat(args[1], 64)
			if valueErr != nil {
				panic(valueErr)
			}
			pm := services.GetPrometheusManager()
			pm.AddMetric(name, value)
		},
	}
}

func removeMetricCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a Prometheus metric",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			pm := services.GetPrometheusManager()
			pm.RemoveMetric(name)
		},
	}
}

func listMetricsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all Prometheus metrics",
		Run: func(cmd *cobra.Command, args []string) {
			pm := services.GetPrometheusManager()
			pm.ListMetrics()
		},
	}
}
