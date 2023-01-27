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
	"gitter/internal/glist"

	"github.com/cwdot/go-stdlib/timediff"
)

var allBranches bool
var showArchived bool
var layoutName string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allBranches, "all", "a", false, "Show all git branches; default branches defined in config")
	listCmd.Flags().BoolVarP(&showArchived, "archived", "", false, "Show archived branches")
	listCmd.Flags().StringVarP(&layoutName, "layout", "", "default", "Customize layout; default is 'default'")
}

func open() (*config.ActiveRepo, *git.Repository, []config.Column, error) {
	if underTest {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, nil, err
		}
		env := filepath.Join(home, ".env")
		return config.OpenCustom(env, layoutName, showArchived)
	}
	return config.OpenDefault(layoutName, showArchived)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		activeRepo, g, layout, err := open()
		if err != nil {
			log.Panic(err)
		}

		rows, err := glist.SortBranches(activeRepo, g, allBranches)
		if err != nil {
			log.Panic(err)
		}

		bench := glist.CreateTable(layout)
		for _, row := range rows {
			branch := row.Branch

			var name, description, links string
			if branch != nil {
				name = branch.Name
				description = branch.Description
				links = glist.GenerateLinks(activeRepo.Repo.BaseLinks, branch.Links)
			}

			rootRow := activeRepo.Repo.RootBranch == name

			var active string
			if row.IsHead {
				active = "*"
			} else if rootRow {
				active = "M"
			} else if row.Archived {
				active = "A"
			} else {
				active = " "
			}

			var commitDateS, relDateS string
			commitDate := row.LastCommit.Committer.When
			if !commitDate.IsZero() {
				commitDateS = commitDate.Format(timediff.TimeLayout)
				now := time.Now()
				relDateS = timediff.Compute(commitDate, now, timediff.EpochRounding())
			}

			output := make(map[config.ColumnKind]string)
			output[config.Active] = active
			output[config.LastHash] = row.Hash
			output[config.LastHashShort] = row.Hash[0:7]
			output[config.Project] = row.Project
			output[config.Name] = name
			output[config.Description] = description
			output[config.LastCommitted] = strings.TrimSpace(row.LastCommit.Message)
			output[config.CommittedDate] = commitDateS
			output[config.RelDate] = relDateS
			output[config.Tracking] = row.TrackingBranch
			output[config.Links] = links

			colors := colorDataRow(row.Project, rootRow, commitDate)
			bench.Append(output, colors)
		}

		bench.Render()
	},
}

// TODO: not aware of different layouts
func colorDataRow(project string, isRoot bool, date time.Time) *glist.RowColor {
	rc := glist.NewRowColor()
	rc.Colors[config.Project] = colorStringByHash(project)
	if isRoot {
		rc.DefaultStyle = tw.Bold
	}
	if date.Before(time.Now().Add(time.Hour * 24 * -30)) {
		rc.Colors[config.CommittedDate] = tw.FgHiRedColor
		rc.Colors[config.RelDate] = tw.FgHiRedColor
	}
	return rc
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
