package cli

import "github.com/spf13/cobra"

func ConnLogzCmds() ([]*cobra.Command, error) {
	prometheuzCmd, prometheuzCmdErr := ConnPrometheusCmd()
	if prometheuzCmdErr != nil {
		return nil, prometheuzCmdErr
	}

	return []*cobra.Command{
		prometheuzCmd,
	}, nil
}

func ConnPrometheusCmd() (*cobra.Command, error) {
	var route string
	var port int

	prometheuzCmd := &cobra.Command{
		Use:     "prometheuz",
		Aliases: []string{"prometheus", "prom", "metrics"},
		Annotations: GetDescriptions([]string{
			"Exposes metrics to Prometheus",
			"Exposes metrics to Prometheus",
		}, false),

		Run: func(cmd *cobra.Command, args []string) {
			//route, _ := cmd.Flags().GetString("route")
			//port, _ := cmd.Flags().GetInt("port")
			//Prometheuz(route, port)
		},
	}
	prometheuzCmd.Flags().StringVarP("route", "r", "", "Route to expose metrics")
	prometheuzCmd.Flags().IntVarP(&clear, "port", "p", 0, "Port to expose metrics")

	return prometheuzCmd, nil
}
