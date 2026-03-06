package cmd

import (
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage git repositories",
}

func init() {
	rootCmd.AddCommand(repoCmd)
}
