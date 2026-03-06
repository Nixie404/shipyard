package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var deployCmd = &cobra.Command{
	Use:   "deploy <stack>",
	Short: "Deploy a stack (shortcut for 'stack deploy')",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg := mustLoadConfig()

		s := cfg.FindStack(name)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found. Run 'yardctl stack list' to see available stacks.\n", name)
			os.Exit(1)
		}

		if err := stack.Deploy(s); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		mustSaveConfig(cfg)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
