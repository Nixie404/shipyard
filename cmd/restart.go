package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var restartCmd = &cobra.Command{
	Use:   "restart <stack> [service]",
	Short: "Restart services in a stack",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		service := ""
		if len(args) > 1 {
			service = args[1]
		}

		cfg := mustLoadConfig()

		s := cfg.FindStack(stackName)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found\n", stackName)
			os.Exit(1)
		}

		if err := stack.Restart(s, service); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
