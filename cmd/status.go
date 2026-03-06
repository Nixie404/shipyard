package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
	"shipyard/internal/stack"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show overall status of repositories and stacks",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		fmt.Printf("📂 Data directory:   %s\n", cfg.DataDir)
		fmt.Printf("📁 Repos directory:  %s\n", cfg.ReposDir)
		fmt.Printf("📋 Config file:      %s\n", config.ConfigPath())
		fmt.Println()

		// Repos
		fmt.Printf("📦 Repositories: %d\n", len(cfg.Repos))
		for _, r := range cfg.Repos {
			fmt.Printf("   • %s (%s)\n", r.Name, r.URL)
		}
		fmt.Println()

		// Stacks
		deployed := 0
		for _, s := range cfg.Stacks {
			if s.Deployed {
				deployed++
			}
		}
		fmt.Printf("🐳 Stacks: %d total, %d deployed\n", len(cfg.Stacks), deployed)

		if len(cfg.Stacks) > 0 {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "   NAME\tSTATUS")
			for _, s := range cfg.Stacks {
				status := "⏹  stopped"
				if s.Deployed {
					status = "▶  running"
				}
				fmt.Fprintf(w, "   %s\t%s\n", s.Name, status)
			}
			w.Flush()
		}
		fmt.Println()

		// Docker status
		fmt.Print("🐋 Docker: ")
		dockerCmd := exec.Command("docker", "info", "--format", "{{.ServerVersion}}")
		out, err := dockerCmd.Output()
		if err != nil {
			fmt.Println("❌ not reachable")
		} else {
			fmt.Printf("✅ v%s", string(out))
		}

		// Docker Compose status
		fmt.Print("📝 Docker Compose: ")
		composeCmd := exec.Command("docker", "compose", "version", "--short")
		out, err = composeCmd.Output()
		if err != nil {
			fmt.Println("❌ not available")
		} else {
			fmt.Printf("✅ v%s", string(out))
		}

		// Git status
		fmt.Print("🔧 Git: ")
		gitCmd := exec.Command("git", "--version")
		out, err = gitCmd.Output()
		if err != nil {
			fmt.Println("❌ not available")
		} else {
			fmt.Printf("✅ %s", string(out))
		}

		// Check for running containers per deployed stack
		if deployed > 0 {
			fmt.Println("\n📊 Running containers:")
			for _, s := range cfg.Stacks {
				if s.Deployed {
					fmt.Printf("\n   [%s]\n", s.Name)
					stack.Ps(&s)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
