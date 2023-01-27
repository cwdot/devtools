package config

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func OpenDefault(showArchived bool) (*Layout, *git.Repository, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, nil, nil
	}
	return OpenCustom(path, showArchived)
}

func OpenCustom(path string, showArchived bool) (*Layout, *git.Repository, error) {
	conf, err := DefaultConfigFile()
	if err != nil {
		return nil, nil, err
	}

	branches := make(map[string]*branchPair)

	repo := getRepo(conf, path)
	sets := []BranchSet{repo.Active}
	if showArchived {
		sets = append(sets, repo.Archived)
	}

	for idx, set := range sets {
		for project, rBranches := range set {
			for _, rBranch := range rBranches {
				if _, ok := branches[rBranch.Name]; ok {
					panic("Branches with the same name: " + rBranch.Name)
				}

				archived := false
				if idx == 1 {
					archived = true
				}

				branches[rBranch.Name] = &branchPair{
					Project:  project,
					Branch:   rBranch,
					Archived: archived,
				}
			}
		}
	}

	gr := Layout{
		Repo:     repo,
		branches: branches,
	}

	cork, err := git.PlainOpen(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open git repo")
	}

	return &gr, cork, nil
}

func getRepo(conf *Config, path string) *Repo {
	repos := make(map[string]*Repo)
	for _, repo := range conf.Repos {
		home := processHome(repo.Home)
		repos[home] = &repo
		if home == path {
			return &repo
		}
	}
	return nil
}

func processHome(value string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "error finding home dir").Error())
	}
	return strings.ReplaceAll(value, "$HOME", home)
}

type branchPair struct {
	Project  string
	Branch   Branch
	Archived bool
}
type Layout struct {
	Repo     *Repo
	branches map[string]*branchPair
}

func (r *Layout) FindBranch(shortName string) (*Branch, string, bool, bool) {
	if val, ok := r.branches[shortName]; ok {
		return &val.Branch, val.Project, val.Archived, ok
	}
	return nil, "", false, false
}
