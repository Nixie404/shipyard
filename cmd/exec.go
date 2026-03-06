package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var execCmd = &cobra.Command{
	Use:   "exec <stack> <service> [command...]",
	Short: "Execute a command in a running service container",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		service := args[1]
		command := []string{"sh"}
		if len(args) > 2 {
			command = args[2:]
		}

		cfg := mustLoadConfig()

		s := cfg.FindStack(stackName)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found\n", stackName)
			os.Exit(1)
		}

		if err := stack.Exec(s, service, command); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
