package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit yardctl configuration",
}

func init() {
	rootCmd.AddCommand(configCmd)
}
