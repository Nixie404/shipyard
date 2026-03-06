package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Show real-time Docker events",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📡 Listening for Docker events (Ctrl+C to stop)...")
		dockerCmd := exec.Command("docker", "events", "--format", "{{.Time}} {{.Type}} {{.Action}} {{.Actor.Attributes.name}}")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr
		if err := dockerCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(eventsCmd)
}
