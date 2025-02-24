package cli

import (
	"github.com/faelmori/logz/internal/services"
	"github.com/spf13/cobra"
)

func ConnLogzCmds() []*cobra.Command {
	return []*cobra.Command{
		ConnPrometheusCmd(),
	}
}

func ConnPrometheusCmd() *cobra.Command {
	var route string
	var port int

	prometheuzCmd := &cobra.Command{
		Use:     "prometheuz",
		Aliases: []string{"prometheus", "prom", "metrics"},
		Annotations: GetDescriptions([]string{
			"Exposes metrics to Prometheus",
			"Exposes metrics to Prometheus",
		}, false),

		RunE: func(cmd *cobra.Command, args []string) error {
			return services.Prometheuz(route, port)
		},
	}
	prometheuzCmd.Flags().StringVarP(&route, "route", "r", "", "Route to expose metrics")
	prometheuzCmd.Flags().IntVarP(&port, "port", "p", 0, "Port to expose metrics")

	return prometheuzCmd
}
