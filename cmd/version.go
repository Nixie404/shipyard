package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"
var Commit = "dev"
var BuildDate = "unknown"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of yardctl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("yardctl (Shipyard) v%s\n", Version)
		fmt.Printf("Commit:  %s\n", Commit)
		fmt.Printf("Built:   %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
