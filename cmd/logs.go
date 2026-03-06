package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/stack"
)

var logsFollow bool

var logsCmd = &cobra.Command{
	Use:   "logs <stack> [service]",
	Short: "Show logs for a stack or specific service",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		service := ""
		if len(args) > 1 {
			service = args[1]
		}

		cfg := mustLoadConfig()

		s := cfg.FindStack(name)
		if s == nil {
			fmt.Fprintf(os.Stderr, "Stack %q not found\n", name)
			os.Exit(1)
		}

		if err := stack.Logs(s, service, logsFollow); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	rootCmd.AddCommand(logsCmd)
}
