package cmd

import (
	"github.com/faelmori/logz/cmd/cli"
	"github.com/spf13/cobra"
)

type Logz struct{}

func (m *Logz) Alias() string {
	return "logs"
}
func (m *Logz) ShortDescription() string {
	return "LoggerLogz and logs manager"
}
func (m *Logz) LongDescription() string {
	return "The \"logz\" command-line interface (CLI) is an intuitive and user-friendly logger and log management module designed for developers. Integrated with Prometheus for monitoring, \"logz\" ensures comprehensive log management and is compatible with other plugins and the Go programming language making it a versatile tool for maintaining system health and performance."
}
func (m *Logz) Usage() string {
	return "logz [command] [args]"
}
func (m *Logz) Examples() []string {
	return []string{"logz show all", "lg error 'error message'"}
}
func (m *Logz) Active() bool {
	return true
}
func (m *Logz) Module() string {
	return "logz"
}
func (m *Logz) Execute() error {
	return m.Command().Execute()
}
func (m *Logz) Command() *cobra.Command {
	var logType, message, name, show, clearLogs, archive string
	var filter []string
	var follow, quiet bool

	cmd := &cobra.Command{
		Use:         m.Module(),
		Annotations: cli.GetDescriptions([]string{m.LongDescription(), m.ShortDescription()}, false),
		Run: func(cmd *cobra.Command, args []string) {
			//return logzCmd.NewLogger([]string{logType, message, name, strconv.FormatBool(quiet), show, strconv.FormatBool(follow), clearLogs, archive}...))
		},
	}

	cmd.Flags().StringVarP(&logType, "type", "t", "", "Tipo de log")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Mensagem de log")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Nome do módulo de log")
	cmd.Flags().StringVarP(&show, "show", "s", "", "Mostrar logs")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Seguir logs")
	cmd.Flags().StringVarP(&clearLogs, "clear", "c", "", "Limpar logs")
	cmd.Flags().StringVarP(&archive, "archive", "z", "", "Arquivar logs")
	cmd.Flags().StringArrayVarP(&filter, "filter", "F", []string{}, "Filtrar logs")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Modo silencioso")

	cmd.AddCommand(cli.LogzCmds()...)
	cmd.AddCommand(cli.ServiceCmd())
	cmd.AddCommand(cli.MetricsCmd())

	setUsageDefinition(cmd)

	for _, c := range cmd.Commands() {
		setUsageDefinition(c)
	}

	return cmd
}
func (m *Logz) concatenateExamples() string {
	examples := ""
	for _, example := range m.Examples() {
		examples += string(example) + "\n  "
	}
	return examples
}

func RegX() *Logz {
	return &Logz{}
}
