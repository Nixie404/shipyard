package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
)

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		fmt.Printf("Config file: %s\n\n", config.ConfigPath())

		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	},
}

func init() {
	configCmd.AddCommand(configViewCmd)
}
