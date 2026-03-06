package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
	"shipyard/internal/privilege"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize yardctl configuration and data directories",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DefaultConfig()
		esc := privilege.Detect()

		fmt.Printf("⚡ Privilege escalation: %s\n\n", esc.Label())

		// Create directories
		dirs := []string{cfg.DataDir, cfg.ReposDir}
		for _, dir := range dirs {
			if err := esc.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
				os.Exit(1)
			}
			fmt.Printf("✅ Created %s\n", dir)
		}

		// Make the data dir owned by the current user so we don't need root for everyday ops
		if !privilege.IsRoot() {
			uid := os.Getuid()
			gid := os.Getgid()
			if err := esc.Chown(cfg.DataDir, uid, gid); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not chown %s: %v\n", cfg.DataDir, err)
			} else {
				fmt.Printf("✅ Owned %s by current user\n", cfg.DataDir)
			}
		}

		// Write config
		cfgPath := config.ConfigPath()
		cfgDir := config.ConfigDir()
		if err := esc.MkdirAll(cfgDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config dir %s: %v\n", cfgDir, err)
			os.Exit(1)
		}

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling config: %v\n", err)
			os.Exit(1)
		}

		if err := esc.WriteFile(cfgPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✅ Config written to %s\n", cfgPath)

		fmt.Println("\n🚀 Shipyard initialized! Next steps:")
		fmt.Println("  yardctl repo add <git-url>")
		fmt.Println("  yardctl repo sync")
		fmt.Println("  yardctl stack list")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
