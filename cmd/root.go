package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"tidyfiles/config"

	"github.com/spf13/cobra"
)

var version = "0.0.1"
var (
	verboseFlag bool
	quietFlag   bool
	configPath  string
)

var rootCmd = &cobra.Command{
	Use:     "tidyfiles",
	Version: version,
	Short:   "File organizer",
	Long:    "A command line tool that organizes files in directories based on configurable rules",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verboseFlag && quietFlag {
			return fmt.Errorf("--verbose and --quiet are mutually exclusive")
		}

		level := slog.LevelInfo
		switch {
		case verboseFlag:
			level = slog.LevelDebug
		case quietFlag:
			level = slog.LevelError
		}
		slog.SetLogLoggerLevel(level)

		_, err := config.LoadConfig(configPath)
		if err != nil {
			return err
		}
		return nil
	},
	Args: cobra.ArbitraryArgs,
	Run:  runCommand,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Debug logs enabled")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Only error logs enabled")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Configuration file used")
}

func runCommand(cmd *cobra.Command, args []string) {
	fmt.Println("RUNNING ROOT CMD")
	for i, j := range args {
		fmt.Println(i, j)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an issue '%s'", err)
		os.Exit(1)
	}
}
