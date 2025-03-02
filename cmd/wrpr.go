package main

import (
	"github.com/faelmori/logz/cmd/cli"
	"github.com/spf13/cobra"
)

// Logz represents the main structure for the logz command-line interface.
type Logz struct{}

// Alias returns the alias for the logz command.
func (m *Logz) Alias() string {
	return "logs"
}

// ShortDescription provides a brief description of the logz command.
func (m *Logz) ShortDescription() string {
	return "LoggerLogz and logs manager"
}

// LongDescription provides a detailed description of the logz command.
func (m *Logz) LongDescription() string {
	return "The \"logz\" command-line interface (CLI) is an intuitive and user-friendly logger and log management module designed for developers. Integrated with Prometheus for monitoring, \"logz\" ensures comprehensive log management and is compatible with other plugins and the Go programming language making it a versatile tool for maintaining system health and performance."
}

// Usage returns the usage information for the logz command.
func (m *Logz) Usage() string {
	return "logz [command] [args]"
}

// Examples returns example usages of the logz command.
func (m *Logz) Examples() []string {
	return []string{"logz show all", "lg error 'error message'"}
}

// Active indicates whether the logz command is active.
func (m *Logz) Active() bool {
	return true
}

// Module returns the module name for the logz command.
func (m *Logz) Module() string {
	return "logz"
}

// Execute runs the logz command.
func (m *Logz) Execute() error {
	return m.Command().Execute()
}

// Command creates and returns the main cobra.Command for the logz CLI.
func (m *Logz) Command() *cobra.Command {
	var logType, message, name, show, clearLogs, archive string
	var filter []string
	var follow, quiet bool

	cmd := &cobra.Command{
		Use:         m.Module(),
		Annotations: cli.GetDescriptions([]string{m.LongDescription(), m.ShortDescription()}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Placeholder for the command execution logic
			// return logzCmd.NewLogger([]string{logType, message, name, strconv.FormatBool(quiet), show, strconv.FormatBool(follow), clearLogs, archive}...))
		},
	}

	// Define flags for the logz command
	cmd.Flags().StringVarP(&logType, "type", "t", "", "Log type")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Log message")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Log module name")
	cmd.Flags().StringVarP(&show, "show", "s", "", "Show logs")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow logs")
	cmd.Flags().StringVarP(&clearLogs, "clear", "c", "", "Clear logs")
	cmd.Flags().StringVarP(&archive, "archive", "z", "", "Archive logs")
	cmd.Flags().StringArrayVarP(&filter, "filter", "F", []string{}, "Filter logs")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	// Add subcommands to the logz command
	cmd.AddCommand(cli.LogzCmds()...)
	cmd.AddCommand(cli.ServiceCmd())
	cmd.AddCommand(cli.MetricsCmd())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(cmd)
	for _, c := range cmd.Commands() {
		setUsageDefinition(c)
	}

	return cmd
}

// concatenateExamples concatenates example usages into a single string.
func (m *Logz) concatenateExamples() string {
	examples := ""
	for _, example := range m.Examples() {
		examples += string(example) + "\n  "
	}
	return examples
}

// RegX returns a new instance of the Logz struct.
func RegX() *Logz {
	return &Logz{}
}
