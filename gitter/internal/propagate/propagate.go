package propagate

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"gitter/internal/config"

	"github.com/cwdot/go-stdlib/wood"
)

func Propagate(activeRepo *config.ActiveRepo, project string, dryRun bool) error {
	branches, ok := activeRepo.FindTree(project)
	if !ok {
		return errors.New("no project found")
	}

	names := make([]string, 0, len(branches))
	for _, item := range branches {
		names = append(names, item.Name)
	}

	wood.Infof("Propagation: %s", strings.Join(names, " => "))
	if dryRun {
		wood.Warnf("Dryrun: true")
	}

	gx := gitX{DryRun: dryRun}

	parent := branches[0].Name
	for _, branch := range branches[1:] {
		err := gx.checkout(branch.Name)
		if err != nil {
			return errors.Wrap(err, "checkout failed")
		}

		err = gx.rebase(branch.Name, parent)
		if err != nil {
			return errors.Wrap(err, "rebase failed")
		}

		parent = branch.Name
	}

	return nil
}

type gitX struct {
	DryRun bool
}

func (x *gitX) checkout(branch string) error {
	err := x.run("checkout", branch)
	if err != nil {
		return errors.Wrap(err, "checkout failed")
	}
	wood.Infof("Checked out: %s", branch)
	return nil
}

func (x *gitX) rebase(branch string, parent string) error {
	err := x.run("rebase", parent)
	if err != nil {
		return errors.Wrap(err, "rebase failed")
	}
	wood.Infof("Pulled `%s` and rebased `%s`", branch, parent)
	return nil
}

func (x *gitX) pullRebase(branch string, parent string) error {
	err := x.run("pull", "--rebase", parent)
	if err != nil {
		return errors.Wrap(err, "pulled --rebase failed")
	}
	wood.Infof("Pulled `%s` and rebased `%s`", branch, parent)
	return nil
}

// Inject replaces placeholders in the template and writes the result into the output file
func (x *gitX) run(args ...string) error {
	if x.DryRun {
		return nil
	}

	cmd := exec.Command("/opt/homebrew/bin/git", args...)

	wood.Debugf("Executing git: %s with %s", cmd, args)

	var outs bytes.Buffer
	var errs bytes.Buffer
	cmd.Stdout = &outs
	cmd.Stderr = &errs
	cmd.Dir = "/Users/indy/.env"

	err := cmd.Run()
	if err != nil {
		os.Stderr.Write(outs.Bytes())
		os.Stderr.Write(errs.Bytes())
		return errors.Wrap(err, "git call failed")
	}

	wood.Debugf("Executed git: %s", outs.String())

	return nil
}
