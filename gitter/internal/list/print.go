package list

import (
	"hash/fnv"
	"strconv"
	"strings"
	"time"

	"github.com/cwdot/go-stdlib/timediff"
	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5"
	tw "github.com/olekukonko/tablewriter"

	"gitter/internal/config"
)

type PrintOpts struct {
	Layout      []config.Column
	AllBranches bool
	NoTrackers  bool
}

func PrintBranches(activeRepo *config.ActiveRepo, g *git.Repository, opts PrintOpts) {
	rows, err := getGitBranchRows(activeRepo, g, opts)
	if err != nil {
		wood.Fatal(err)
	}

	bench := createTable(opts.Layout)
	for _, row := range rows {
		branch := row.Branch

		var name, description, links string
		if branch != nil {
			name = branch.Name
			description = branch.Description
			links = GenerateLinks(&activeRepo.Repo.BaseLinks, branch)
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
		output[config.RootDrift] = strconv.Itoa(row.RootDrift)
		output[config.RootDriftDesc] = row.RootDriftDesc
		output[config.RootTracking] = row.RootTracking
		output[config.RemoteDrift] = strconv.Itoa(row.RemoteDrift)
		output[config.RemoteDriftDesc] = row.RemoteDriftDesc
		output[config.RemoteTracking] = row.RemoteTracking
		output[config.Links] = links

		colors := colorDataRow(row.Project, rootRow, commitDate, row.RootDrift, row.RemoteDrift)
		bench.Append(output, colors)
	}

	bench.Render()
}

func colorDataRow(project string, isRoot bool, date time.Time, rootDrift int, remoteDrift int) *RowColor {
	rc := NewRowColor()
	rc.Colors[config.Project] = colorStringByHash(project)
	if isRoot {
		rc.DefaultStyle = tw.Bold
	}
	if date.Before(time.Now().Add(time.Hour * 24 * -30)) {
		rc.Colors[config.CommittedDate] = tw.FgHiRedColor
		rc.Colors[config.RelDate] = tw.FgHiRedColor
	}
	if rootDrift > 5 {
		rc.Colors[config.RootDrift] = tw.FgHiYellowColor
		rc.Colors[config.RootDriftDesc] = tw.FgHiYellowColor
	} else if rootDrift < 5 {
		rc.Colors[config.RootDrift] = tw.FgHiRedColor
		rc.Colors[config.RootDriftDesc] = tw.FgHiRedColor
	}
	if remoteDrift > 5 {
		rc.Colors[config.RemoteDrift] = tw.FgHiYellowColor
		rc.Colors[config.RemoteDriftDesc] = tw.FgHiYellowColor
	} else if remoteDrift < 5 {
		rc.Colors[config.RemoteDrift] = tw.FgHiRedColor
		rc.Colors[config.RemoteDriftDesc] = tw.FgHiRedColor
	}
	return rc
}

func colorStringByHash(text string) int {
	h := fnv.New32a()
	h.Write([]byte(text))
	color := h.Sum32() % 5
	return tw.FgRedColor + int(color)
}
