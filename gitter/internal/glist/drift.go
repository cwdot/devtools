package glist

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/pkg/errors"
)

func computeDrift(g *git.Repository, tracking string, childTop *object.Commit) (int, string, error) {
	trackingHash, err := g.ResolveRevision(plumbing.Revision(tracking))
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed to resolve drift revision: %s", "shortName")
	}

	return findHashInCommits(childTop, *trackingHash)
}

func findHashInCommits(commit *object.Commit, target plumbing.Hash) (int, string, error) {
	count := 0
	found := false
	iter := object.NewCommitPreorderIter(commit, nil, nil)
	err := iter.ForEach(func(comm *object.Commit) error {
		if comm.Hash != target {
			count++
			return nil
		}
		found = true
		return storer.ErrStop
	})
	if err != nil {
		return 0, "", errors.Wrapf(err, "failed to iterate commit: %s", commit.Message)
	}

	switch {
	case !found:
		return 0, "no overlap", nil
	case count > 0:
		return count, fmt.Sprintf("%d ahead", count), nil
	case count < 0:
		return count, fmt.Sprintf("%d behind", count), nil
	default:
		return count, "same", nil
	}
}
