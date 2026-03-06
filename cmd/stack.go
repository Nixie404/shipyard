package cmd

import (
	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage Docker Compose stacks",
}

func init() {
	rootCmd.AddCommand(stackCmd)
}
