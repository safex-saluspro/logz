package cli

import (
	lgzCmd "github.com/faelmori/logz/internal/cmd"
	lgzUtl "github.com/faelmori/logz/internal/utils"
	"github.com/spf13/cobra"
)

func LogzCmds() ([]*cobra.Command, error) {
	return []*cobra.Command{
		logzDebugCmd(),
		logzWriterCmd(),
	}, nil
}

func logzDebugCmd() *cobra.Command {
	var message, module string
	var quiet bool

	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Log a debug message",
		Long:  "Log a debug message",
		RunE: func(cmd *cobra.Command, args []string) error {
			lgzUtl.

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
		Use:   "writer",
		Short: "Log IO writer",
		Long:  "Log IO writer to the log module",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = lgzCmd
			return nil
		},
	}

	cmd.Flags().StringVarP(&message, "message", "m", "", "Message to log")
	cmd.Flags().StringVarP(&module, "context", "n", "logz", "Context of the log module to link the log to")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")

	return cmd
}
