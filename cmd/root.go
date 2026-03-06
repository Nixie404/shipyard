package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "yardctl",
	Short: "Docker-centered automated deployment from git repositories",
	Long: `yardctl (Shipyard) is a docker-centered automated deployment tool
that pulls configuration from git repositories (including private repos),
discovers Docker Compose stacks, and manages their lifecycle.

Typical workflow:
  yardctl init
  yardctl repo add git@github.com:org/deploy.git
  yardctl repo sync
  yardctl stack list
  yardctl stack deploy api
  yardctl logs api
  yardctl status`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// mustLoadConfig loads config or exits with error.
func mustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

// mustSaveConfig saves config or exits with error.
func mustSaveConfig(cfg *config.Config) {
	if err := cfg.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}
}
