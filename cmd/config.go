package cmd

import (
	"log/slog"
	"tidyfiles/config"

	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config [init|show]",
	Short: "Configuration related tools",
	Args:  cobra.ExactArgs(1),
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default config file",
	Args:  cobra.MaximumNArgs(1),
	Run:   runInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows current configuration",
	Run:   runShow,
}

func runInit(cmd *cobra.Command, args []string) {
	var err error
	if len(args) == 1 {
		err = config.AppConfig.SaveToFile(args[0])
	} else {
		err = config.AppConfig.SaveToFile("")
	}
	if err != nil {
		slog.Error("Could not create configuration file", "err", err)
	}
}

func runShow(cmd *cobra.Command, args []string) {
	config.AppConfig.PrintConfiguration()
}

func init() {
	configCommand.AddCommand(configInitCmd)
	configCommand.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCommand)
}
