package gitp

import (
	"fmt"
	"strings"

	"github.com/cwdot/go-stdlib/wood"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

func computeDrift(g *git.Repository, master, feature string) (int, string, error) {
	wood.Infof("Computing drift between: %s and %s", master, feature)

	// behind is master; ahead is feature branch
	behind, ahead, err := calc(g, master, feature)
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

func calcMergeBase(g *git.Repository, masterBranch *plumbing.Reference, featureBranch *plumbing.Reference) (plumbing.Hash, error) {
	mHash := masterBranch.Hash()
	fHash := featureBranch.Hash()

	masterCommit, err := g.CommitObject(mHash)
	if err != nil {
		return plumbing.ZeroHash, errors.Wrapf(err, "failed to find master commit: %s", mHash)
	}

	featureCommit, err := g.CommitObject(fHash)
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

func calc(g *git.Repository, master, feature string) (int, int, error) {
	masterBranch, err := g.Reference(plumbing.NewBranchReferenceName(master), true) //"refs/heads/master"
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find master branch: %s", master)
	}

	featureBranch, err := g.Reference(plumbing.NewBranchReferenceName(feature), true) //"refs/heads/feature"
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find feature branch: %s", feature)
	}

	baseCommit, err := calcMergeBase(g, masterBranch, featureBranch)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find merge base: %s", baseCommit)
	}

	masterCommits, err := getBranchCommits(g, master)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find commits ahead: %s", master)
	}

	featureCommits, err := getBranchCommits(g, feature)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to find commits ahead: %s", feature)
	}

	masterCommitsAheadCount := 0
	featureCommitsAheadCount := 0

	for _, commit := range masterCommits {
		if commit.Hash == baseCommit {
			break
		}
		masterCommitsAheadCount++
	}

	for _, commit := range featureCommits {
		if commit.Hash == baseCommit {
			break
		}
		featureCommitsAheadCount++
	}

	return masterCommitsAheadCount, featureCommitsAheadCount, nil
}

func getBranchCommits(g *git.Repository, branchName string) ([]*object.Commit, error) {
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
	commits := make([]*object.Commit, 0, 100)
	err = history.ForEach(func(hc *object.Commit) error {
		commits = append(commits, hc)
		return nil
	})
	return commits, nil
}
