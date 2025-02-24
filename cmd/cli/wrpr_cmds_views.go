package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func ViewersCmds() []*cobra.Command {
	return []*cobra.Command{
		analyzeLogzCmd(),
		viewLogzUiCmd(),
		showLogzCmd(),
	}, nil
}

func showLogzCmd() *cobra.Command {
	showCmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"list", "view"},
		Short:   "Show logs",
		Long:    "Show logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			nameFlagValue, _ := cmd.Flags().GetString("name")
			followFlagValue, _ := cmd.Flags().GetBool("follow")
			linesFlagValue, _ := cmd.Flags().GetInt("lines")
			sinceFlagValue, _ := cmd.Flags().GetString("since")
			untilFlagValue, _ := cmd.Flags().GetString("until")
			colorsFlagValue, _ := cmd.Flags().GetBool("colors")
			newArgs := []string{nameFlagValue, fmt.Sprintf("%t", followFlagValue), fmt.Sprintf("%d", linesFlagValue), sinceFlagValue, untilFlagValue, fmt.Sprintf("%t", colorsFlagValue)}
			args = append(args, newArgs...)
			_, showLogzErr := ShowLog(args...)
			return showLogzErr
		},
	}

	showCmd.Flags().StringP("name", "n", "all", "Log name")
	showCmd.Flags().BoolP("follow", "f", false, "Follow logs")
	showCmd.Flags().IntP("lines", "l", 10, "Number of lines to show")
	showCmd.Flags().StringP("since", "s", "", "Show logs since a specific time")
	showCmd.Flags().StringP("until", "u", "", "Show logs until a specific time")
	showCmd.Flags().BoolP("colors", "c", true, "Show logs with colors")

	return showCmd, nil
}

func analyzeLogzCmd() *cobra.Command {

	analyzeCmd := &cobra.Command{
		Use:     "analyze",
		Aliases: []string{"report"},
		Short:   "Analyze logs",
		Long:    "Analyze logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("file flag is required")
			}

			return analyzeLog(file)
		},
	}
	analyzeCmd.Flags().StringP("file", "f", "", "Log file to analyze")

	return analyzeCmd, nil
}

func viewLogzUiCmd() *cobra.Command {
	viewCmd := &cobra.Command{
		Use:     "ui",
		Aliases: []string{"web", "interface"},
		Short:   "View logs in a web interface",
		Long:    "View logs in a web interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			//viewLogzUi()
			return fmt.Errorf("in development")
		},
	}

	return viewCmd, nil
}
