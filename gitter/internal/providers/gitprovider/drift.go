package gitprovider

import (
	"fmt"
	"strings"

	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

func NewDrifter(g *git.Repository) *DriftCalculator {
	return &DriftCalculator{
		g:             g,
		branchCommits: make(map[string][]plumbing.Hash),
	}
}

type DriftCalculator struct {
	g             *git.Repository
	branchCommits map[string][]plumbing.Hash
}

func (c *DriftCalculator) computeDrift(master string, feature string) (int, string, error) {
	wood.Debugf("Computing drift between: %s and %s", master, feature)

	// behind is master; ahead is feature branch
	behind, ahead, err := c.calc(master, feature)
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed to resolve drift revision: %s and %s", master, feature)
	}
	var buf strings.Builder

	if ahead > 0 {
		buf.WriteString(fmt.Sprintf("%d ahead", ahead))
	}
	if behind > 0 {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d behind", behind))
	}
	return ahead - behind, buf.String(), nil
}

func (c *DriftCalculator) calcMergeBase(masterBranch *plumbing.Reference, featureBranch *plumbing.Reference) (plumbing.Hash, error) {
	mHash := masterBranch.Hash()
	fHash := featureBranch.Hash()

	masterCommit, err := c.g.CommitObject(mHash)
	if err != nil {
		return plumbing.ZeroHash, errors.Wrapf(err, "failed to find master commit: %s", mHash)
	}

	featureCommit, err := c.g.CommitObject(fHash)
	if err != nil {
		return plumbing.ZeroHash, errors.Wrapf(err, "failed to find feature commit: %s", fHash)
	}

	// Find the merge base between the two branches
	commits, err := masterCommit.MergeBase(featureCommit)
	if err != nil {
		return plumbing.ZeroHash, errors.Wrapf(err, "failed to find merge base between %s and %s", mHash, fHash)
	}

	if len(commits) == 1 {
		return commits[0].Hash, nil
	}
	return plumbing.ZeroHash, errors.Errorf("unexpected number of commits: %d", len(commits))
}

func (c *DriftCalculator) calc(master string, feature string) (int, int, error) {
	masterBranch, err := c.g.Reference(plumbing.NewBranchReferenceName(master), true)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find master branch: %s", master)
	}

	featureBranch, err := c.g.Reference(plumbing.NewBranchReferenceName(feature), true)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find feature branch: %s", feature)
	}

	baseCommit, err := c.calcMergeBase(masterBranch, featureBranch)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find merge base: %s", baseCommit)
	}

	masterCommits, err := c.getBranchCache(master)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find commits ahead: %s", master)
	}

	featureCommits, err := c.getBranchCache(feature)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find commits ahead: %s", feature)
	}

	masterCommitsAheadCount := 0
	featureCommitsAheadCount := 0

	for _, commit := range masterCommits {
		if commit == baseCommit {
			break
		}
		masterCommitsAheadCount++
	}

	for _, commit := range featureCommits {
		if commit == baseCommit {
			break
		}
		featureCommitsAheadCount++
	}

	return masterCommitsAheadCount, featureCommitsAheadCount, nil
}

func (c *DriftCalculator) getBranchCache(branch string) ([]plumbing.Hash, error) {
	if cache, ok := c.branchCommits[branch]; ok {
		return cache, nil
	}

	commits, err := getBranchCommits(c.g, branch)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find commits ahead: %s", branch)
	}
	c.branchCommits[branch] = commits

	return commits, nil
}

func getBranchCommits(g *git.Repository, branchName string) ([]plumbing.Hash, error) {
	refName := plumbing.NewBranchReferenceName(branchName)

	// Resolve the reference to a commit
	ref, err := g.Reference(refName, true)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find reference: %s", refName)
	}

	// Get the commit object for the branch's head
	commitObj, err := g.CommitObject(ref.Hash())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find commit object for branch: %s", refName)
	}

	// Iterate through and list the commits in the branch

	history := object.NewCommitPreorderIter(commitObj, nil, nil)
	commits := make([]plumbing.Hash, 0, 100)
	err = history.ForEach(func(hc *object.Commit) error {
		commits = append(commits, hc.Hash)
		return nil
	})
	return commits, nil
}
