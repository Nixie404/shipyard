package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var psCmd = &cobra.Command{
	Use:   "ps [stack]",
	Short: "List running containers, optionally for a specific stack",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		if len(args) == 1 {
			s := cfg.FindStack(args[0])
			if s == nil {
				fmt.Fprintf(os.Stderr, "Stack %q not found\n", args[0])
				os.Exit(1)
			}
			stack.Ps(s)
			return
		}

		// Show all deployed stacks
		for _, s := range cfg.Stacks {
			if s.Deployed {
				fmt.Printf("━━━ %s ━━━\n", s.Name)
				stack.Ps(&s)
				fmt.Println()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
