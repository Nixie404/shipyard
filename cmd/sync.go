package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/repo"
	"shipyard/internal/stack"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync all repos and rediscover stacks (shortcut for 'repo sync')",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		fmt.Println("🔄 Syncing all repositories...")
		if err := repo.SyncAll(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

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
	rootCmd.AddCommand(syncCmd)
}
