package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
	"shipyard/internal/privilege"
)

var purgeForce bool

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove all yardctl data, config, and systemd units (undo init)",
	Long: `Completely removes everything created by 'yardctl init':
  - /etc/yardctl/          (configuration)
  - /var/lib/yardctl/      (cloned repos, stack state)
  - systemd timer/service  (if installed)

This does NOT remove running Docker containers from deployed stacks.
Run 'yardctl stack destroy <stack>' first if you want to tear those down.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()
		esc := privilege.Detect()

		// Check for deployed stacks
		deployed := 0
		for _, s := range cfg.Stacks {
			if s.Deployed {
				deployed++
			}
		}

		if deployed > 0 {
			fmt.Printf("⚠  %d stack(s) are currently marked as deployed:\n", deployed)
			for _, s := range cfg.Stacks {
				if s.Deployed {
					fmt.Printf("   • %s\n", s.Name)
				}
			}
			fmt.Println("\n   Their containers will keep running but yardctl will no longer manage them.")
			fmt.Println("   Run 'yardctl stack destroy <stack>' first to tear them down cleanly.")
			fmt.Println()
		}

		// Confirm unless --force
		if !purgeForce {
			fmt.Println("This will permanently remove:")
			fmt.Printf("   📁 %s\n", config.ConfigDir())
			fmt.Printf("   📁 %s\n", config.DataDir())
			fmt.Println("   ⏱  yardctl-sync.timer/service (if installed)")
			fmt.Println()
			fmt.Print("Are you sure? Type 'yes' to confirm: ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer != "yes" {
				fmt.Println("Aborted.")
				return
			}
			fmt.Println()
		}

		fmt.Printf("⚡ Privilege escalation: %s\n\n", esc.Label())

		// Stop and disable systemd timer if present
		fmt.Print("⏱  Disabling systemd timer... ")
		_ = esc.Run("systemctl", "stop", "yardctl-sync.timer")
		_ = esc.Run("systemctl", "disable", "yardctl-sync.timer")
		_ = esc.Run("systemctl", "stop", "yardctl-sync.service")
		fmt.Println("done")

		// Remove data directory
		dataDir := config.DataDir()
		fmt.Printf("🗑  Removing %s... ", dataDir)
		if err := esc.Run("rm", "-rf", dataDir); err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Println("done")
		}

		// Remove config directory
		cfgDir := config.ConfigDir()
		fmt.Printf("🗑  Removing %s... ", cfgDir)
		if err := esc.Run("rm", "-rf", cfgDir); err != nil {
			fmt.Printf("❌ %v\n", err)
		} else {
			fmt.Println("done")
		}

		fmt.Println()
		fmt.Println("✅ Purge complete. All yardctl state has been removed.")
		fmt.Println("   Run 'yardctl init' to start fresh.")
	},
}

func init() {
	purgeCmd.Flags().BoolVarP(&purgeForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(purgeCmd)
}
