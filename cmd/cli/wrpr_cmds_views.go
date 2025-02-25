package cli

import (
	"fmt"
	logzCmd "github.com/faelmori/logz/internal/cmd"
	"github.com/spf13/cobra"
)

func ViewersCmds() []*cobra.Command {
	return []*cobra.Command{
		analyzeLogzCmd(),
		viewLogzUiCmd(),
		showLogzCmd(),
	}
}

func showLogzCmd() *cobra.Command {
	var nameFlagValue, sinceFlagValue, untilFlagValue string
	var followFlagValue, colorsFlagValue bool
	var linesFlagValue int

	showCmd := &cobra.Command{
		Use:         "show",
		Aliases:     []string{"list", "view"},
		Annotations: GetDescriptions([]string{"Show logs", "Show logs"}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			newArgs := []string{nameFlagValue, fmt.Sprintf("%t", followFlagValue), fmt.Sprintf("%d", linesFlagValue), sinceFlagValue, untilFlagValue, fmt.Sprintf("%t", colorsFlagValue)}
			args = append(args, newArgs...)
			_, showLogzErr := logzCmd.NewLogz().ShowLog(args...)
			return showLogzErr
		},
	}

	showCmd.Flags().StringVarP(&nameFlagValue, "name", "n", "all", "Log name")
	showCmd.Flags().BoolVarP(&followFlagValue, "follow", "f", false, "Follow logs")
	showCmd.Flags().IntVarP(&linesFlagValue, "lines", "l", 10, "Number of lines to show")
	showCmd.Flags().StringVarP(&sinceFlagValue, "since", "s", "", "Show logs since a specific time")
	showCmd.Flags().StringVarP(&untilFlagValue, "until", "u", "", "Show logs until a specific time")
	showCmd.Flags().BoolVarP(&colorsFlagValue, "colors", "c", true, "Show logs with colors")

	return showCmd
}

func analyzeLogzCmd() *cobra.Command {
	var file string

	analyzeCmd := &cobra.Command{
		Use:         "analyze",
		Aliases:     []string{"report"},
		Annotations: GetDescriptions([]string{"Analyze logs", "Analyze logs"}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("file flag is required")
			}
			return logzCmd.NewLogz().AnalyzeLog(file)
		},
	}
	analyzeCmd.Flags().StringVarP(&file, "file", "f", "", "Log file to analyze")

	return analyzeCmd
}

func viewLogzUiCmd() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:         "ui",
		Aliases:     []string{"web", "interface"},
		Annotations: GetDescriptions([]string{"View logs in a web interface", "View logs in a web interface"}, false),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("in development")
		},
	}

	return viewCmd
}
