package list

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"gitter/internal/config"
)

type gitBranchMetadata struct {
	Branch          *config.Branch
	Project         string
	Archived        bool
	IsHead          bool
	Hash            string
	LastCommit      *object.Commit
	RootTracking    string
	RootDrift       int
	RootDriftDesc   string
	RemoteTracking  string
	RemoteDrift     int
	RemoteDriftDesc string
}

func getGitBranchRows(layout *config.ActiveRepo, g *git.Repository, allBranches bool) ([]*gitBranchMetadata, error) {
	iter, err := g.Branches()
	if err != nil {
		wood.Fatal(err)
	}

	head, err := g.Head()
	if err != nil {
		wood.Fatal(err)
	}

	refs := make([]*gitBranchMetadata, 0, 20)

	// child is main branch we're on
	// parent/root is usually master
	// remote is usually child's remote branch
	err = iter.ForEach(func(r *plumbing.Reference) error {
		shortName := r.Name().Short()
		branch, project, archived, ok := layout.FindBranch(shortName)
		if !allBranches && !ok {
			return nil
		}

		lastChildCommit, err := g.CommitObject(r.Hash())
		if err != nil {
			return errors.Wrapf(err, "failed to find commit object for last commit: %s", shortName)
		}

		var rootTracking, rootDriftDesc, remoteTracking, remoteDriftDesc string
		var rootDrift, remoteDrift int

		if layout.Repo.RootBranch != "" && layout.Repo.RootBranch != shortName {
			rootTracking = layout.Repo.RootBranch
			rootDrift, rootDriftDesc, err = computeDrift(g, layout.Repo.RootBranch, lastChildCommit)
			if err != nil {
				return errors.Wrapf(err, "failed to find drift for root: %s", shortName)
			}
		}

		if branch.RemoteBranch != "" {
			remoteTracking = branch.RemoteBranch
			remoteDrift, remoteDriftDesc, err = computeDrift(g, branch.RemoteBranch, lastChildCommit)
			if err != nil {
				return errors.Wrapf(err, "failed to find drift for remote: %s", shortName)
			}
		}

		refs = append(refs, &gitBranchMetadata{
			Branch:          branch,
			Project:         project,
			Archived:        archived,
			LastCommit:      lastChildCommit,
			RootTracking:    rootTracking,
			RootDrift:       rootDrift,
			RootDriftDesc:   rootDriftDesc,
			RemoteTracking:  remoteTracking,
			RemoteDrift:     remoteDrift,
			RemoteDriftDesc: remoteDriftDesc,
			Hash:            r.Hash().String(),
			IsHead:          r.Hash() == head.Hash(),
		})
		return nil
	})
	if err != nil {
		wood.Fatal(err)
	}

	sortBranches(layout.Repo.RootBranch, refs)

	return refs, nil
}

func sortBranches(rootBranchName string, refs []*gitBranchMetadata) {
	sort.Slice(refs, func(i, j int) bool {
		// handle missing parts
		if refs[i].Branch == nil && refs[j].Branch == nil {
			return refs[i].Hash < refs[j].Hash
		} else if refs[i].Branch != nil && refs[j].Branch == nil {
			return true
		} else if refs[i].Branch == nil && refs[j].Branch != nil {
			return false
		}

		// we have all parts; compare in proper sequence
		switch {
		case refs[i].Branch.Name == rootBranchName:
			return true
		case refs[j].Branch.Name == rootBranchName:
			return false
		case refs[i].Project != refs[j].Project:
			return refs[i].Project < refs[j].Project
		case refs[i].Branch.Name != refs[j].Branch.Name:
			return refs[i].Branch.Name < refs[j].Branch.Name
		default:
			return false
		}
	})
}

func GenerateLinks(base config.BaseLinks, links config.BranchLinks) string {
	if links.Pr != "" {
		return createCsvLinks(base.PrBase, links.Pr)
	}
	if links.Jira != "" {
		return createCsvLinks(base.JiraBase, links.Jira)
	}
	return ""
}

func createCsvLinks(base string, csvLinks string) string {
	newLinks := make([]string, 0, 5)
	items := strings.Split(csvLinks, ",")
	for _, item := range items {
		newLinks = append(newLinks, fmt.Sprintf("%s/%s", base, item))
	}
	return strings.Join(newLinks, " ")
}
