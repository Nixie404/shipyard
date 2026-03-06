package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system dependencies and configuration health",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🩺 Shipyard Doctor")
		fmt.Println("==================")
		allGood := true

		// Check Docker
		fmt.Print("Docker.............. ")
		if out, err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}").Output(); err != nil {
			fmt.Println("❌ not reachable")
			allGood = false
		} else {
			fmt.Printf("✅ v%s", string(out))
		}

		// Check Docker Compose
		fmt.Print("Docker Compose...... ")
		if out, err := exec.Command("docker", "compose", "version", "--short").Output(); err != nil {
			fmt.Println("❌ not available")
			allGood = false
		} else {
			fmt.Printf("✅ v%s", string(out))
		}

		// Check Git
		fmt.Print("Git................. ")
		if out, err := exec.Command("git", "version").Output(); err != nil {
			fmt.Println("❌ not available")
			allGood = false
		} else {
			fmt.Printf("✅ %s", string(out))
		}

		// Check config dir
		fmt.Print("Config directory.... ")
		cfgPath := config.ConfigPath()
		if _, err := os.Stat(cfgPath); err != nil {
			fmt.Printf("❌ %s not found (run 'yardctl init')\n", cfgPath)
			allGood = false
		} else {
			fmt.Printf("✅ %s\n", cfgPath)
		}

		// Check data dir
		fmt.Print("Data directory...... ")
		dataDir := config.DataDir()
		if _, err := os.Stat(dataDir); err != nil {
			fmt.Printf("❌ %s not found (run 'yardctl init')\n", dataDir)
			allGood = false
		} else {
			fmt.Printf("✅ %s\n", dataDir)
		}

		// Check repos dir
		fmt.Print("Repos directory..... ")
		reposDir := config.ReposDir()
		if _, err := os.Stat(reposDir); err != nil {
			fmt.Printf("❌ %s not found (run 'yardctl init')\n", reposDir)
			allGood = false
		} else {
			fmt.Printf("✅ %s\n", reposDir)
		}

		// Check SSH agent
		fmt.Print("SSH Agent........... ")
		if os.Getenv("SSH_AUTH_SOCK") != "" {
			fmt.Println("✅ running")
		} else {
			fmt.Println("⚠  not detected (private repos may fail)")
		}

		// Check systemd timer
		fmt.Print("Sync timer.......... ")
		if out, err := exec.Command("systemctl", "is-active", "yardctl-sync.timer").Output(); err != nil {
			fmt.Println("⏹  not installed")
		} else {
			fmt.Printf("✅ %s", string(out))
		}

		fmt.Println()
		if allGood {
			fmt.Println("✅ All checks passed!")
		} else {
			fmt.Println("⚠  Some checks failed. Address the issues above.")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
