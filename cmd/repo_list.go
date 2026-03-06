package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var repoListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List registered repositories",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := mustLoadConfig()

		if len(cfg.Repos) == 0 {
			fmt.Println("No repositories registered. Use 'yardctl repo add <url>' to add one.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tURL\tBRANCH")
		fmt.Fprintln(w, "----\t---\t------")
		for _, r := range cfg.Repos {
			branch := r.Branch
			if branch == "" {
				branch = "main"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.URL, branch)
		}
		w.Flush()
	},
}

func init() {
	repoCmd.AddCommand(repoListCmd)
}
