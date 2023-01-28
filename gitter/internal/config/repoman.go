package config

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

func OpenDefault(layoutName string, showArchived bool) (*ActiveRepo, *git.Repository, []Column, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, nil
	}
	return OpenCustom(path, layoutName, showArchived)
}

func OpenCustom(path string, layoutName string, showArchived bool) (*ActiveRepo, *git.Repository, []Column, error) {
	conf, err := DefaultConfigFile()
	if err != nil {
		return nil, nil, nil, err
	}

	var layout []Column
	if layoutName == "default" {
		layout = DefaultLayout()
	} else {
		var ok bool
		layout, ok = conf.Layouts[layoutName]
		if !ok {
			return nil, nil, nil, errors.Errorf("Cannot find layout in config: %s", layoutName)
		}
	}

	repo, err := getRepo(conf, path)
	if err != nil {
		return nil, nil, nil, err
	}

	sets := []BranchSet{repo.Active}
	if showArchived {
		sets = append(sets, repo.Archived)
	}

	branches := make(map[string]*branchPair)
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

	gr := ActiveRepo{
		Repo:     repo,
		branches: branches,
	}

	cork, err := git.PlainOpen(path)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to open git repo")
	}

	return &gr, cork, layout, nil
}

func getRepo(conf *Config, path string) (*Repo, error) {
	repos := make(map[string]*Repo)
	for _, repo := range conf.Repos {
		home := processHome(repo.Home)
		repos[home] = &repo
		if home == path {
			return &repo, nil
		}
	}
	return nil, errors.Errorf("failed to find repo: %s", path)
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
type ActiveRepo struct {
	Repo     *Repo
	branches map[string]*branchPair
}

func (r *ActiveRepo) FindBranch(shortName string) (*Branch, string, bool, bool) {
	if val, ok := r.branches[shortName]; ok {
		return &val.Branch, val.Project, val.Archived, ok
	}
	return nil, "", false, false
}
