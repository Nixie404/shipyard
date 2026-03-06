package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
)

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open the configuration file in your editor",
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}

		path := config.ConfigPath()
		fmt.Printf("Opening %s with %s...\n", path, editor)

		editCmd := exec.Command(editor, path)
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr

		if err := editCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	configCmd.AddCommand(configEditCmd)
}
