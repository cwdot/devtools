package config

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"

	"github.com/cwdot/stdlib-go/wood"
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

	repo, realPath, err := getRepo(conf, path)
	if err != nil {
		return nil, nil, nil, err
	}

	sets := []BranchSet{repo.Active}
	if showArchived {
		sets = append(sets, repo.Archived)
	}

	branches := make(map[string]*BranchRef)
	projects := make(map[string][]Branch)
	for idx, set := range sets {
		for project, rBranches := range set {
			projects[project] = rBranches
			for _, rBranch := range rBranches {
				if _, ok := branches[rBranch.Name]; ok {
					panic("Branches with the same name: " + rBranch.Name)
				}

				archived := false
				if idx == 1 {
					archived = true
				}

				branches[rBranch.Name] = &BranchRef{
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
		projects: projects,
		trees:    repo.Trees,
	}

	cork, err := git.PlainOpen(realPath)
	if err != nil {
		return nil, nil, nil, errors.Wrapf(err, "failed to open git repo: %s", realPath)
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

	return &gr, cork, layout, nil
}

func getRepo(conf *Config, path string) (*Repo, string, error) {
	candidates, err := computeCandidates(path)
	if err != nil {
		return nil, "", err
	}

	names := make([]string, 0, len(conf.Repos))
	repos := make(map[string]*Repo)
	for _, repo := range conf.Repos {
		home := repo.Home
		repos[home] = &repo
		for _, candidate := range candidates {
			if home == candidate {
				wood.Debugf("Found repo: %v", candidate)
				return &repo, candidate, nil
			}
		}
		names = append(names, repo.Name)
	}

	wood.Infof("Path candidates: %v", candidates)
	wood.Infof("Repos: %v", names)
	return nil, "", errors.Errorf("failed to find repo")
}

func computeCandidates(path string) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Errorf("failed to find home: %s", homeDir)
	}
	candidates := make([]string, 0, 5)
	for {
		if path == homeDir || path == "/" {
			break // never home dir
		}
		candidates = append(candidates, path)
		path = filepath.Dir(path)
	}
	return candidates, nil
}

type BranchRef struct {
	Project  string
	Branch   Branch
	Archived bool
}
type ActiveRepo struct {
	Repo     *Repo
	branches map[string]*BranchRef
	projects BranchSet
	trees    TreeSet
}

func (r *ActiveRepo) FindBranch(shortName string) (*BranchRef, bool) {
	if val, ok := r.branches[shortName]; ok {
		return val, ok
	}
	return nil, false
}

func (r *ActiveRepo) FindByProject(projectName string) ([]Branch, bool) {
	if val, ok := r.projects[projectName]; ok {
		return val, ok
	}
	return nil, false
}

func (r *ActiveRepo) FindTree(projectName string) ([]TreeBranch, bool) {
	if val, ok := r.trees[projectName]; ok {
		return val, ok
	}
	return nil, false
}
