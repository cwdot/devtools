package config

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

func OpenDefault() (*Layout, *git.Repository, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, nil, nil
	}
	return OpenCustom(path)
}

func OpenCustom(path string) (*Layout, *git.Repository, error) {
	conf, err := DefaultConfigFile()
	if err != nil {
		return nil, nil, err
	}

	repo := getRepo(conf, path)

	cork, err := git.PlainOpen(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open git repo")
	}

	branches := make(map[string]*branchPair)

	sets := []BranchSet{repo.Active, repo.Archived}
	for _, set := range sets {
		for project, rBranches := range set {
			for _, rBranch := range rBranches {
				branches[rBranch.Name] = &branchPair{
					Project: project,
					Branch:  rBranch,
				}
			}
		}
	}

	gr := Layout{
		Repo:     repo,
		branches: branches,
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
	Project string
	Branch  Branch
}
type Layout struct {
	Repo     *Repo
	branches map[string]*branchPair
}

func (r *Layout) FindBranch(refName plumbing.ReferenceName) (*Branch, string, bool) {
	if val, ok := r.branches[refName.Short()]; ok {
		return &val.Branch, val.Project, ok
	}
	return nil, "", false
}
