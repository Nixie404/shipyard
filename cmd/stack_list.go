package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var stackListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List discovered stacks",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		if len(cfg.Stacks) == 0 {
			fmt.Println("No stacks discovered. Run 'yardctl repo sync' to discover stacks.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tREPO\tSTATUS\tCOMPOSE FILE")
		fmt.Fprintln(w, "----\t----\t------\t------------")
		for _, s := range cfg.Stacks {
			status := "stopped"
			if s.Deployed {
				status = "deployed"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.RepoName, status, s.ComposePath)
		}
		w.Flush()
	},
}

func init() {
	stackCmd.AddCommand(stackListCmd)
}
