package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var stackInfoCmd = &cobra.Command{
	Use:   "info <stack>",
	Short: "Show information about a stack",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg := mustLoadConfig()

		s := cfg.FindStack(name)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found\n", name)
			os.Exit(1)
		}

		status := "stopped"
		if s.Deployed {
			status = "deployed"
		}

		fmt.Printf("Stack:        %s\n", s.Name)
		fmt.Printf("Repository:   %s\n", s.RepoName)
		fmt.Printf("Compose File: %s\n", s.ComposePath)
		fmt.Printf("Status:       %s\n", status)
	},
}

func init() {
	stackCmd.AddCommand(stackInfoCmd)
}
