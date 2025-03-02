package cli

import (
	"fmt"
	"github.com/faelmori/logz/internal/logger"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// LogzCmds returns the CLI commands for different log levels and management.
func LogzCmds() []*cobra.Command {
	return []*cobra.Command{
		newLogCmd("debug", []string{"dbg"}),
		newLogCmd("info", []string{"inf"}),
		newLogCmd("warn", []string{"wrn"}),
		newLogCmd("error", []string{"err"}),
		newLogCmd("fatal", []string{"ftl"}),
		watchLogsCmd(),
		startServiceCmd(),
		stopServiceCmd(),
		rotateLogsCmd(),
		checkLogSizeCmd(),
		archiveLogsCmd(),
	}
}

// newLogCmd configures a command for a specific log level.
func newLogCmd(level string, aliases []string) *cobra.Command {
	var metaData, ctx map[string]string
	var msg string

	cmd := &cobra.Command{
		Use:     level,
		Aliases: aliases,
		Annotations: GetDescriptions(
			[]string{"Logs a " + level + " level message"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Error initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				return
			}

			logr := logger.NewLogger(config)
			for k, v := range metaData {
				logr.SetMetadata(k, v)
			}
			ctxInterface := make(map[string]interface{})
			for k, v := range ctx {
				ctxInterface[k] = v
			}
			switch level {
			case "debug":
				logr.Debug(msg, ctxInterface)
			case "info":
				logr.Info(msg, ctxInterface)
			case "warn":
				logr.Warn(msg, ctxInterface)
			case "error":
				logr.Error(msg, ctxInterface)
			case "fatal":
				logr.FatalC(msg, ctxInterface)
			default:
				logr.Info(msg, ctxInterface)
			}
		},
	}

	cmd.Flags().StringToStringVarP(&metaData, "metadata", "m", nil, "Metadata to include")
	cmd.Flags().StringToStringVarP(&ctx, "context", "c", nil, "Context for the log")
	cmd.Flags().StringVarP(&msg, "msg", "M", "", "Log message")

	return cmd
}

// rotateLogsCmd allows manual log rotation.
func rotateLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "rotate",
		Annotations: GetDescriptions(
			[]string{"Rotates logs that exceed the configured size"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Error initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				return
			}

			err = logger.CheckLogSize(config)
			if err != nil {
				fmt.Printf("Error rotating logs: %v\n", err)
			} else {
				fmt.Println("Logs rotated successfully!")
			}
		},
	}
}

// checkLogSizeCmd checks the current log size.
func checkLogSizeCmd() *cobra.Command {
	return &cobra.Command{
		Use: "check-size",
		Annotations: GetDescriptions(
			[]string{"Checks the log size without taking any action"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Error initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				return
			}

			logDir := config.DefaultLogPath()
			logSize, err := logger.GetLogDirectorySize(filepath.Dir(logDir)) // Add this function to logger
			if err != nil {
				fmt.Printf("Error calculating log size: %v\n", err)
				return
			}

			fmt.Printf("The total log size in directory '%s' is: %d bytes\n", filepath.Dir(logDir), logSize)
		},
	}
}

// archiveLogsCmd allows manual log archiving.
func archiveLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "archive",
		Annotations: GetDescriptions(
			[]string{"Manually archives all logs"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			err := logger.ArchiveLogs(nil)
			if err != nil {
				fmt.Printf("Error archiving logs: %v\n", err)
			} else {
				fmt.Println("Logs archived successfully!")
			}
		},
	}
}

// watchLogsCmd monitors logs in real-time.
func watchLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "watch",
		Aliases: []string{"w"},
		Annotations: GetDescriptions(
			[]string{"Monitors logs in real-time"},
			false,
		),
		Run: func(cmd *cobra.Command, args []string) {
			configManager := logger.NewConfigManager()
			if configManager == nil {
				fmt.Println("Error initializing ConfigManager.")
				return
			}
			cfgMgr := *configManager

			config, err := cfgMgr.LoadConfig()
			if err != nil {
				fmt.Printf("Error loading configuration: %v\n", err)
				return
			}

			logFilePath := config.DefaultLogPath()
			reader := logger.NewFileLogReader()
			stopChan := make(chan struct{})

			// Capture signals for interruption
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChan
				close(stopChan)
			}()

			fmt.Println("Monitoring started (Ctrl+C to exit):")
			if err := reader.Tail(logFilePath, stopChan); err != nil {
				fmt.Printf("Error monitoring logs: %v\n", err)
			}

			// Wait a small delay to finish
			time.Sleep(500 * time.Millisecond)
		},
	}
}
