package cli

import "github.com/spf13/cobra"

func LogzCmd() (*cobra.Command, error) {
	var logType, message, logModuleName string
	var show, follow, clearLogs, archive, quiet bool

	cmd := &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias(), "log", "lg", "l"},
		Example: m.concatenateExamples(),
		Short:   m.ShortDescription(),
		Long:    m.LongDescription(),
		RunE: func(cmd *cobra.Command, args []string) error {

		},
	}

	// Define as flags do comando diretamente
	cmd.Flags().StringVarP(&logType, "type", "t", "", "Tipo de log")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Mensagem de log")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Nome do módulo de log")
	cmd.Flags().StringVarP(&show, "show", "s", "", "Mostrar logs")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Seguir logs")
	cmd.Flags().StringVarP("clear", "c", "", "Limpar logs")
	cmd.Flags().StringVarP("archive", "z", "", "Arquivar logs")
	cmd.Flags().StringVarP("filter", "l", "", "Filtrar logs")
	cmd.Flags().BoolVarP("quiet", "q", false, "Modo silencioso")

	return cmd, nil
}
