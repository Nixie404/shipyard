package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/repo"
)

var repoRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Short:   "Remove a registered repository",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		cfg := mustLoadConfig()

		if !cfg.RemoveRepo(name) {
			fmt.Fprintf(os.Stderr, "Repository %q not found\n", name)
			os.Exit(1)
		}

		// Remove directory from disk
		if err := repo.RemoveDir(cfg, name); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to remove repo directory: %v\n", err)
		}

		mustSaveConfig(cfg)
		fmt.Printf("✅ Repository %q removed\n", name)
	},
}

func init() {
	repoCmd.AddCommand(repoRemoveCmd)
}
