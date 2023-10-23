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
	"gitter/internal/jirap"
	"gitter/internal/providers/gitprovider"
	"gitter/internal/providers/jiraprovider"
)

func PrintBranches(activeRepo *config.ActiveRepo, g *git.Repository, opts config.PrintOpts) {
	rows, err := gitprovider.GetGitBranchRows(activeRepo, g, opts)
	if err != nil {
		wood.Fatal(err)
	}

	jiras := make([]string, 0, 10)
	for _, row := range rows {
		if row.BranchConf.Jira != "" {
			jiras = append(jiras, row.BranchConf.Jira)
		} else if key := jirap.SafeExtract(activeRepo.Repo.Jira, row.BranchName); key != "" {
			jiras = append(jiras, key)
		}
	}
	issues, err := jiraprovider.GetIssues(activeRepo.Repo.Jira, jiras...)
	if err != nil {
		wood.Fatal(err)
	}

	bench := createTable(opts.Layout)
	for _, row := range rows {
		branchName := row.BranchName

		var description, links, jiraStatus string
		branchConf := row.BranchConf
		if branchConf.Name != "" {
			description = branchConf.Description
			links = gitprovider.GenerateLinks(activeRepo.Repo, branchConf)
			if issue, ok := issues[branchConf.Jira]; ok {
				jiraStatus = issue.Fields.Status.Name
			}
		}

		rootRow := activeRepo.Repo.RootBranch == branchName

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
		output[config.Name] = branchName
		output[config.Description] = description
		output[config.LastCommitted] = strings.TrimSpace(row.LastCommit.Message)
		output[config.CommittedDate] = commitDateS
		output[config.RelDate] = relDateS
		output[config.RootDrift] = strconv.Itoa(row.RootDrift)
		output[config.MainDriftDesc] = row.RootDriftDesc
		output[config.MainTracking] = row.RootTracking
		output[config.RemoteDrift] = strconv.Itoa(row.RemoteDrift)
		output[config.RemoteDriftDesc] = row.RemoteDriftDesc
		output[config.RemoteTracking] = row.RemoteTracking
		output[config.JiraStatus] = jiraStatus
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
		rc.Colors[config.MainDriftDesc] = tw.FgHiYellowColor
	} else if rootDrift < 5 {
		rc.Colors[config.RootDrift] = tw.FgHiRedColor
		rc.Colors[config.MainDriftDesc] = tw.FgHiRedColor
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
