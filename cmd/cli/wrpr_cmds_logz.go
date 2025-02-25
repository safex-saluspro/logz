package cli

import (
	"github.com/spf13/cobra"
)

func LogzCmds() []*cobra.Command {
	return []*cobra.Command{
		logzDebugCmd(),
		logzWriterCmd(),
	}
}

func logzDebugCmd() *cobra.Command {
	var message, module string
	var quiet bool

	cmd := &cobra.Command{
		Use:         "debug",
		Aliases:     []string{"dbg"},
		Annotations: GetDescriptions([]string{"Logs a debug message", "Logs a debug message"}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Message to log")
	cmd.Flags().StringVarP(&module, "context", "n", "logz", "Context of the log module to link the log to")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	return cmd
}

func logzWriterCmd() *cobra.Command {
	var message, module string
	var quiet bool

	cmd := &cobra.Command{
		Use:         "writer",
		Aliases:     []string{"wrt"},
		Annotations: GetDescriptions([]string{"Logs a writer message", "Logs a writer message"}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Message to log")
	cmd.Flags().StringVarP(&module, "context", "n", "logz", "Context of the log module to link the log to")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	return cmd
}
