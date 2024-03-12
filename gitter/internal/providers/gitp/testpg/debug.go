package testpg

import (
	"fmt"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

func PrintRefs(r *git.Repository) {
	refs, err := r.References()
	if err != nil {
		panic(err)
	}
	printRefs(refs, "refs")
}

func PrintBranches(r *git.Repository) {
	refs, err := r.Branches()
	if err != nil {
		panic(err)
	}
	printRefs(refs, "branches")
}

func printRefs(refs storer.ReferenceIter, t string) {
	count := 0
	err := refs.ForEach(func(r *plumbing.Reference) error {
		fmt.Println(t, r.Name().String())
		count++
		return nil
	})
	if err != nil {
		panic(err)
	}
	if count == 0 {
		wood.Infof("Reference iterator has ZERO %s", t)
	}
}

func PrintCommits(r *git.Repository) {
	refs, err := r.CommitObjects()
	if err != nil {
		panic(err)
	}

	count := 0
	err = refs.ForEach(func(commit *object.Commit) error {
		fmt.Println("commit", commit.Hash.String()[0:7], commit.Message)
		count++
		return nil
	})
	if err != nil {
		panic(err)
	}
	if count == 0 {
		wood.Infof("Reference iterator has ZERO commits")
	}
}
