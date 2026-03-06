package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contextListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List available contexts",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("* local (default)")
		fmt.Println("\nContexts are defined in the configuration. Use 'yardctl config edit' to add remote contexts.")
	},
}

func init() {
	contextCmd.AddCommand(contextListCmd)
}
