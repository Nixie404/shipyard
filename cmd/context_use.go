package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var contextUseCmd = &cobra.Command{
	Use:   "use <cluster>",
	Short: "Switch to a different deployment context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if name == "local" {
			fmt.Printf("✅ Switched to context %q\n", name)
		} else {
			fmt.Printf("⚠ Context %q not found. Only 'local' is currently supported.\n", name)
		}
	},
}

func init() {
	contextCmd.AddCommand(contextUseCmd)
}
