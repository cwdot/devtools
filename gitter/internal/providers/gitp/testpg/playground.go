package testpg

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
)

func New() (*Playground, error) {
	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return nil, fmt.Errorf("error creating the repository: %v", err)
	}

	// Set master as HEAD
	branch := plumbing.NewBranchReferenceName("master")
	err = r.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, branch))
	if err != nil {
		return nil, errors.Wrapf(err, "error setting master as HEAD")
	}

	return &Playground{FS: fs, R: r}, nil
}

type Playground struct {
	FS billy.Filesystem
	R  *git.Repository
}

func (p *Playground) Checkout(branchName string, create bool) error {
	worktree, err := p.R.Worktree()
	if err != nil {
		return errors.Wrapf(err, "error getting worktree for branch: %s", branchName)
	}
	err = worktree.Checkout(&git.CheckoutOptions{
		Create: create,
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
	if err != nil {
		return errors.Wrapf(err, "error checking out branch: %s", branchName)
	}
	return nil
}

func (p *Playground) Master() error {
	branchName := "master"
	worktree, err := p.R.Worktree()
	if err != nil {
		return errors.Wrapf(err, "error getting worktree for branch: %s", branchName)
	}
	// checkout if it exists (init doesn't create master)
	if _, err := p.R.Branch("master"); err == nil {
		err = worktree.Checkout(&git.CheckoutOptions{})
		if err != nil {
			return errors.Wrapf(err, "error checking out branch: %s", branchName)
		}
	}
	return nil
}

func (p *Playground) Commit() error {
	commitTarget := "master"
	if h, err := p.R.Head(); err == nil {
		commitTarget = h.Name().Short()
	}

	worktree, err := p.R.Worktree()
	if err != nil {
		return errors.Wrapf(err, "error getting worktree")
	}

	status, err := worktree.Status()
	if err != nil {
		return errors.Wrap(err, "error committing work tree")
	}

	files := make([]string, 0, 10)
	for filename, s := range status {
		if s.Staging == git.Unmodified && s.Worktree == git.Unmodified {
			continue
		}
		if s.Staging == git.Renamed {
			files = append(files, fmt.Sprintf("%s -> %s", filename, s.Extra))
		}
		files = append(files, fmt.Sprintf("[%c%c]%s", s.Staging, s.Worktree, filename))
	}

	commitMessage := fmt.Sprintf("Commit %s => %s", commitTarget, strings.Join(files, " "))
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "error committing work tree")
	}
	return nil
}

func (p *Playground) WriteFile(filename string, content string) error {
	f, err := p.FS.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "error opening file: %s", filename)
	}
	_, err = io.WriteString(f, content)
	if err != nil {
		return errors.Wrapf(err, "error writing file: %s", filename)
	}
	if err := f.Close(); err != nil {
		return errors.Wrapf(err, "error closing file: %s", filename)
	}
	return nil
}

func (p *Playground) AddFile(filename string) error {
	worktree, err := p.R.Worktree()
	if err != nil {
		return errors.Wrapf(err, "error getting worktree")
	}

	if _, err := worktree.Add(filename); err != nil {
		return errors.Wrapf(err, "error adding file to work tree: %s", filename)
	}
	return nil
}

func (p *Playground) AddTestFile(filename string) error {
	content := randomString()
	if err := p.WriteFile(filename, content); err != nil {
		return errors.Wrapf(err, "error writing file: %s", filename)
	}
	if err := p.AddFile(filename); err != nil {
		return errors.Wrapf(err, "error adding file: %s", filename)
	}
	return p.Commit()
}
