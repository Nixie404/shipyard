package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/repo"
	"shipyard/internal/stack"
)

var repoSyncCmd = &cobra.Command{
	Use:   "sync [name]",
	Short: "Pull latest changes from repositories and rediscover stacks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		if len(args) == 1 {
			// Sync a specific repo
			name := args[0]
			entry := cfg.FindRepo(name)
			if entry == nil {
				fmt.Fprintf(os.Stderr, "Repository %q not found\n", name)
				os.Exit(1)
			}
			fmt.Printf("🔄 Syncing %s...\n", name)
			if err := repo.SyncRepo(cfg, entry); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Sync all repos
			fmt.Println("🔄 Syncing all repositories...")
			if err := repo.SyncAll(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		}

		// Re-discover stacks
		fmt.Println("🔍 Discovering stacks...")
		discovered, err := stack.Discover(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: stack discovery failed: %v\n", err)
		} else {
			stack.UpdateStacks(cfg, discovered)
			fmt.Printf("   Found %d stack(s)\n", len(cfg.Stacks))
		}

		mustSaveConfig(cfg)
		fmt.Println("✅ Sync complete")
	},
}

func init() {
	repoCmd.AddCommand(repoSyncCmd)
}
