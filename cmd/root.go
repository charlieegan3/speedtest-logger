package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "speedtest-logger",
	Short: "A CLI tool to run speedtest.net tests and post the results to an HTTP endpoint",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
