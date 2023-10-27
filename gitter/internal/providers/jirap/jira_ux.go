package jirap

import (
	"slices"

	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"gitter/internal/config"
	"gitter/internal/util"
)

type Row struct {
	Key    string
	Branch string
	Title  string
	Status string
	Link   string
}

func Build(g *git.Repository, j *config.JiraConfig) []*Row {
	branches := make(map[string]string)
	keys := make([]string, 0, 50)
	iter, err := g.Branches()
	if err != nil {
		wood.Fatal(err)
	}
	err = iter.ForEach(func(r *plumbing.Reference) error {
		shortName := r.Name().Short()
		key := Extract(j.Extraction, shortName)
		if key != "" {
			keys = append(keys, key)
			branches[key] = shortName
		} else {
			branches[key] = ""
		}
		return nil
	})
	if len(keys) == 0 {
		wood.Fatal("No issues found")
	}

	issues, err := GetIssues(j, keys...)
	if err != nil {
		wood.Fatal(err)
	}
	if len(issues) == 0 {
		wood.Fatalf("JIRA returned 0 issues: %s", keys)
	}

	rows := make([]*Row, 0, len(branches))
	for branch, key := range branches {
		var title, status, link string
		if issue, ok := issues[key]; ok {
			title = issue.Fields.Summary
			status = issue.Fields.Status.Name
			link = util.CreateCsvLinks(j.Base, key)
		}
		rows = append(rows, &Row{
			Key:    key,
			Branch: branch,
			Title:  title,
			Status: status,
			Link:   link,
		})
	}

	slices.SortFunc(rows, func(a, b *Row) int {
		if a.Status == b.Status {
			if a.Key < b.Key {
				return -1
			}
			return 0
		}
		if a.Status < b.Status {
			return -1
		}
		return 1
	})

	return rows
}
