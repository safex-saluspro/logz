package cmd

import (
	"github.com/faelmori/logz/cmd/cli"
	logzCmd "github.com/faelmori/logz/internal/cmd"
	"github.com/spf13/cobra"
	"strconv"
)

type Logz struct{}

func RegX() *Logz {
	return &Logz{}
}

func (m *Logz) Alias() string {
	return "logs"
}

func (m *Logz) ShortDescription() string {
	return "Logger and logs manager"
}

func (m *Logz) LongDescription() string {
	return "Logger and logs manager module. It allows to log messages and manage logs easily."
}

func (m *Logz) Usage() string {
	return "logz [command] [args]"
}

func (m *Logz) Examples() []string {
	return []string{"logz show all", "kbx lg error 'error message'"}
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

func (m *Logz) concatenateExamples() string {
	examples := ""
	for _, example := range m.Examples() {
		examples += string(example) + "\n  "
	}
	return examples
}

func (m *Logz) Command() *cobra.Command {
	var logType, message, name, show, clearLogs, archive string
	var filter []string
	var follow, quiet bool

	cmd := &cobra.Command{
		Use:         m.Module(),
		Aliases:     []string{m.Alias(), "log", "lg", "l"},
		Example:     m.concatenateExamples(),
		Annotations: cli.GetDescriptions([]string{m.LongDescription(), m.ShortDescription()}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return logzCmd.NewLogz().Log([]string{logType, message, name, strconv.FormatBool(quiet), show, strconv.FormatBool(follow), clearLogs, archive}...)
		},
	}

	// Define as flags do comando diretamente
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
	cmd.AddCommand(cli.ViewersCmds()...)

	return cmd
}
