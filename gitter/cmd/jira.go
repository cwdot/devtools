package cmd

import (
	"os"
	"slices"

	"github.com/andygrunwald/go-jira"
	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5/plumbing"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"gitter/internal/jirap"

	"gitter/internal/providers/jirap"
)

func init() {
	rootCmd.AddCommand(jiraCmd)
}

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "jira",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		ar, g, _, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		j := ar.Repo.Jira
		if j == nil {
			wood.Fatal("JIRA no configured")
		}

		branches := make(map[string]string)
		keys := make([]string, 0, 50)
		if len(args) > 0 {
			for _, arg := range args {
				key := jirap.Extract(j.Extraction, arg)
				if key != "" {
					keys = append(keys, key)
					branches[key] = arg
				}
			}
		} else {
			iter, err := g.Branches()
			if err != nil {
				wood.Fatal(err)
			}
			err = iter.ForEach(func(r *plumbing.Reference) error {
				shortName := r.Name().Short()
				key := jirap.Extract(j.Extraction, shortName)
				if key != "" {
					keys = append(keys, key)
					branches[key] = shortName
				}
				return nil
			})
		}
		if len(keys) == 0 {
			wood.Fatal("No issues found")
		}

		issues, err := jirap.GetIssuesSlice(j, keys...)
		if err != nil {
			wood.Fatal(err)
		}
		if len(issues) == 0 {
			wood.Fatalf("JIRA returned 0 issues: %s", keys)
		}

		slices.SortFunc(issues, func(a, b jira.Issue) int {
			if a.Fields.Status.Name == b.Fields.Status.Name {
				if a.Key < b.Key {
					return -1
				}
				return 0
			}
			if a.Fields.Status.Name < b.Fields.Status.Name {
				return -1
			}
			return 1
		})

		table := tw.NewWriter(os.Stdout)
		table.SetHeader([]string{"JIRA", "Branch", "Title", "Status"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)
		for _, issue := range issues {
			branch := ""
			if b, ok := branches[issue.Key]; ok {
				branch = b
			}
			table.Append([]string{
				color.Cyan.It(issue.Key),
				color.Yellow.It(branch),
				issue.Fields.Summary,
				color.Magenta.It(issue.Fields.Status.Name),
			})
		}
		table.Render()
	},
}
