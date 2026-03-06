package cmd

import (
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage deployment contexts (clusters)",
}

func init() {
	rootCmd.AddCommand(contextCmd)
}
