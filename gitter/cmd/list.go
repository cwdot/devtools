package cmd

import (
	"hash/fnv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	tw "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gitter/internal/config"
	"gitter/internal/listbranches"
	"gitter/util"
)

var (
	names  = []string{"*", "HASH", "PROJECT", "NAME", "DESCRIPTION", "LAST COMMITTED", "COMMITTED DATE", "REL DATE", "TRACKING", "LINKS"}
	widths = []int{
		1,  // active
		7,  // hash
		20, // project
		30, // name
		30, // description
		30, // last commit
		20, // commit date
		8,  // rel date
		14, // tracking
		30, // links
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	configCmd.AddCommand(configPrintCmd)
}

func open() (*config.Layout, *git.Repository, error) {
	if underTest {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, err
		}
		env := filepath.Join(home, ".env")
		return config.OpenCustom(env)
	}
	return config.OpenDefault()
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		layout, g, err := open()
		if err != nil {
			log.Panic(err)
		}

		rows, err := listbranches.SortBranches(layout, g)
		if err != nil {
			log.Panic(err)
		}

		var data []tableRow

		for _, row := range rows {
			var active string
			if row.IsHead {
				active = "*"
			} else {
				active = " "
			}
			branch := row.Branch

			var name, description, links string
			if branch != nil {
				name = branch.Name
				description = branch.Description
				links = listbranches.GenerateLinks(layout.Repo.BaseLinks, branch.Links)
			}

			var commitDateS, relDateS string
			commitDate := row.LastCommit.Committer.When
			if !commitDate.IsZero() {
				commitDateS = commitDate.Format(util.TimeLayout)
				now := time.Now()
				relDateS = util.TimeDiff(now, commitDate)
			}

			output := []string{
				active,        // active
				row.Hash[0:7], // hash
				row.Project,   // project
				name,          // name
				description,   // description
				strings.TrimSpace(row.LastCommit.Message), // description
				commitDateS,        // commit date
				relDateS,           // rel date
				row.TrackingBranch, // tracking
				links,              // links
			}
			data = append(data, tableRow{
				Data:   output,
				Colors: colorDataRow(row.Project),
			})
		}

		table := tw.NewWriter(os.Stdout)
		table.SetHeader(names)
		for idx, val := range widths {
			table.SetColMinWidth(idx, val)
		}
		table.SetHeaderColor(headerColors()...)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)

		for _, row := range data {
			//if row. {
			//	table.Rich(colorData1, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor}, tablewriter.Colors{}})
			//	table.Rich(colorData2, []tablewriter.Colors{tablewriter.Colors{tablewriter.Normal, tablewriter.FgMagentaColor}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.BgRedColor}, tablewriter.Colors{tablewriter.FgHiGreenColor, tablewriter.Italic, tablewriter.BgHiCyanColor}})
			//}
			if row.Colors == nil {
				table.Append(row.Data)
			} else {
				table.Rich(row.Data, row.Colors)
			}
		}

		table.Render()
	},
}

func headerColors() []tw.Colors {
	return []tw.Colors{
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
		{tw.Bold, tw.FgHiBlueColor},
	}
}

func colorDataRow(project string) []tw.Colors {
	projectColor := colorStringByHash(project)
	return []tw.Colors{
		{tw.Bold, tw.FgHiGreenColor},
		{tw.Normal, tw.FgHiCyanColor},
		{tw.Normal, projectColor},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
		{tw.Normal, tw.Normal},
	}
}

func colorStringByHash(text string) int {
	h := fnv.New32a()
	h.Write([]byte(text))
	color := h.Sum32() % 5
	return tw.FgRedColor + int(color)
}

type tableRow struct {
	Data   []string
	Colors []tw.Colors
}
