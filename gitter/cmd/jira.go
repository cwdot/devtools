package cmd

import (
	"os"

	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/wood"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"gitter/internal/datatable"
	"gitter/internal/providers/jirap"
)

func init() {
	rootCmd.AddCommand(jiraCmd)
	jiraCmd.Flags().BoolP("all", "a", false, "Show all git branches")
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

		all := mustRet(cmd.Flags().GetBool("all"))

		rows := jirap.Build(g, j)
		statusM := datatable.NewMarker()
		statusM.Set("Backlog", color.Cyan)
		statusM.Set("Development", color.Yellow)
		statusM.Set("InReview", color.Cyan)
		statusM.Set("Done", color.Green)

		table := tw.NewWriter(os.Stdout)
		table.SetHeader([]string{"Branch", "JIRA", "Title", "Status", "Link"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)
		for _, row := range rows {
			// filter out branches without JIRAs
			if row.Key == "" && !all {
				continue
			}
			table.Append([]string{
				color.Yellow.It(row.Branch),
				color.Cyan.It(row.Key),
				row.Title,
				statusM.Mark(row.Status),
				row.Link,
			})
		}
		table.Render()
	},
}
