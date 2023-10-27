package cmd

import (
	"os"

	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/wood"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

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

		rows := jirap.Build(g, j)

		table := tw.NewWriter(os.Stdout)
		table.SetHeader([]string{"JIRA", "Branch", "Title", "Status"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)
		for _, row := range rows {
			table.Append([]string{
				color.Cyan.It(row.Key),
				color.Yellow.It(row.Branch),
				row.Title,
				color.Magenta.It(row.Status),
				row.Link,
			})
		}
		table.Render()
	},
}
