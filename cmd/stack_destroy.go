package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var stackDestroyCmd = &cobra.Command{
	Use:   "destroy <stack>",
	Short: "Destroy a stack using docker compose down",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg := mustLoadConfig()

		s := cfg.FindStack(name)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found\n", name)
			os.Exit(1)
		}

		if err := stack.Destroy(s); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		mustSaveConfig(cfg)
	},
}

func init() {
	stackCmd.AddCommand(stackDestroyCmd)
}
