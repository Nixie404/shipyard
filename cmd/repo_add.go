package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/repo"
	"shipyard/internal/stack"
)

var repoAddBranch string

var repoAddCmd = &cobra.Command{
	Use:   "add <git-url>",
	Short: "Add and clone a git repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		cfg := mustLoadConfig()

		fmt.Printf("📦 Adding repository %s...\n", url)
		entry, err := repo.Add(cfg, url, repoAddBranch)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Auto-discover stacks after adding
		fmt.Println("🔍 Discovering stacks...")
		discovered, err := stack.Discover(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: stack discovery failed: %v\n", err)
		} else {
			stack.UpdateStacks(cfg, discovered)
			fmt.Printf("   Found %d stack(s)\n", len(cfg.Stacks))
		}

		mustSaveConfig(cfg)
		fmt.Printf("✅ Repository %q added (branch: %s)\n", entry.Name, entry.Branch)
	},
}

func init() {
	repoAddCmd.Flags().StringVarP(&repoAddBranch, "branch", "b", "main", "Branch to track")
	repoCmd.AddCommand(repoAddCmd)
}
