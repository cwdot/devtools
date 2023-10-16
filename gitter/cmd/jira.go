package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/andygrunwald/go-jira"
	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"

	tw "github.com/olekukonko/tablewriter"

	"gitter/internal/config"
	"gitter/internal/jirap"
	"gitter/internal/providers/jiraprovider"
)

func init() {
	rootCmd.AddCommand(jiraCmd)
	jiraCmd.Flags().BoolP("branches", "b", false, "Extract JIRAs from branches")
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

		var keys []string
		fromBranches, err := cmd.Flags().GetBool("branches")
		if err != nil {
			wood.Fatal(err)
		}
		if fromBranches {
			iter, err := g.Branches()
			if err != nil {
				wood.Fatal(err)
			}
			err = iter.ForEach(func(r *plumbing.Reference) error {
				shortName := r.Name().Short()
				key := jirap.Extract(j.ExtractionExpr, shortName)
				if key != "" {
					keys = append(keys, key)
					branches[key] = shortName
				}
				return nil
			})
		} else {
			inputReader := cmd.InOrStdin()
			if len(args) > 0 && args[0] != "-" {
				file, err := os.Open(args[0])
				if err != nil {
					panic(fmt.Errorf("failed open file: %v", err))
				}
				inputReader = file
			}
			keys = processKeys(j, inputReader)
		}

		issues, err := jiraprovider.GetIssuesSlice(j, keys...)
		if err != nil {
			wood.Fatal(err)
		}
		if len(issues) == 0 {
			wood.Fatal("No issues found")
		}

		slices.SortFunc(issues, func(a, b jira.Issue) int {
			if a.Fields.Status.Name < b.Fields.Status.Name {
				return -1
			}
			if a.Key < b.Key {
				return -1
			}
			return 0
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

func processKeys(j *config.JiraConfig, f io.Reader) []string {
	keys := make([]string, 0, 10)
	input := bufio.NewScanner(f)
	for input.Scan() {
		txt := input.Text()
		if txt == "" {
			break
		}
		key := jirap.Extract(j.ExtractionExpr, txt)
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}
