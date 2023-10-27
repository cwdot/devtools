package gitp

import (
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"

	"gitter/internal/config"
	"gitter/internal/util"

	"github.com/cwdot/go-stdlib/wood"
)

type GitBranchMetadata struct {
	BranchConf      config.Branch
	BranchName      string
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

func GetGitBranchRows(layout *config.ActiveRepo, g *git.Repository, printOpts config.PrintOpts) ([]*GitBranchMetadata, error) {
	iter, err := g.Branches()
	if err != nil {
		wood.Fatal(err)
	}

	head, err := g.Head()
	if err != nil {
		wood.Fatal(err)
	}

	refs := make([]*GitBranchMetadata, 0, 20)

	// child is main branch we're on
	// parent/root is usually master
	// remote is usually child's remote branch
	err = iter.ForEach(func(r *plumbing.Reference) error {
		shortName := r.Name().Short()
		ref, ok := layout.FindBranch(shortName)
		if !printOpts.AllBranches && !ok {
			return nil
		}

		var remoteBranch string
		if ref != nil {
			remoteBranch = ref.Branch.RemoteBranch
		}

		lastChildCommit, err := g.CommitObject(r.Hash())
		if err != nil {
			return errors.Wrapf(err, "failed to find commit object for last commit: %s", shortName)
		}

		var rootTracking, rootDriftDesc, remoteTracking, remoteDriftDesc string
		var rootDrift, remoteDrift int

		lastTime := lastChildCommit.Committer.When
		cutoff := time.Now().Add(-14 * 24 * time.Hour)

		// If we're too old, then just don't try.
		if lastTime.Before(cutoff) || printOpts.NoTrackers {
			rootDriftDesc = "~~"
			remoteDriftDesc = "~~"
		} else {
			if layout.Repo.RootBranch != "" && layout.Repo.RootBranch != shortName {
				rootTracking = layout.Repo.RootBranch
				rootDrift, rootDriftDesc, err = computeDrift(g, layout.Repo.RootBranch, shortName)
				if err != nil {
					wood.Warnf("Failed to find drift for root: %s => %s", shortName, err)
				}
			}

			if remoteBranch != "" {
				remoteTracking = remoteBranch
				remoteDrift, remoteDriftDesc, err = computeDrift(g, remoteBranch, shortName)
				if err != nil {
					wood.Warnf("Failed to find drift for remote: %s => %s", shortName, err)
				}
			}
		}

		bm := &GitBranchMetadata{
			BranchName:      shortName,
			LastCommit:      lastChildCommit,
			RootTracking:    rootTracking,
			RootDrift:       rootDrift,
			RootDriftDesc:   rootDriftDesc,
			RemoteTracking:  remoteTracking,
			RemoteDrift:     remoteDrift,
			RemoteDriftDesc: remoteDriftDesc,
			Hash:            r.Hash().String(),
			IsHead:          r.Hash() == head.Hash(),
		}
		if ref != nil {
			if ref.Branch.Name != "" {
				bm.BranchConf = ref.Branch
			}
			bm.Project = ref.Project
			bm.Archived = ref.Archived
		}
		refs = append(refs, bm)
		return nil
	})
	if err != nil {
		wood.Fatal(err)
	}

	sortBranches(layout.Repo.RootBranch, refs)

	return refs, nil
}

func sortBranches(rootBranchName string, refs []*GitBranchMetadata) {
	sort.Slice(refs, func(i, j int) bool {
		bnI := refs[i].BranchConf.Name
		bnJ := refs[j].BranchConf.Name

		// handle missing parts
		if bnI == "" && bnJ == "" {
			return refs[i].Hash < refs[j].Hash
		} else if bnI != "" && bnJ == "" {
			return true
		} else if bnI == "" && bnJ != "" {
			return false
		}

		// we have all parts; compare in proper sequence
		switch {
		case bnI == rootBranchName:
			return true
		case bnJ == rootBranchName:
			return false
		case refs[i].Project != refs[j].Project:
			return refs[i].Project < refs[j].Project
		case bnI != bnJ:
			return bnI < bnJ
		default:
			return false
		}
	})
}

func GenerateLinks(base *config.Repo, links config.Branch) string {
	if links.Pr != "" {
		return util.CreateCsvLinks(base.BaseLinks.PrBase, links.Pr)
	}
	if links.Jira != "" {
		if base.Jira == nil {
			return "config err"
		}
		return util.CreateCsvLinks(base.Jira.Base, links.Jira)
	}
	return ""
}
