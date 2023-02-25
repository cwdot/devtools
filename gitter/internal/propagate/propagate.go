package propagate

import (
	"strings"

	"github.com/pkg/errors"
	"gitter/internal/config"

	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/proc"
	"github.com/cwdot/go-stdlib/wood"
)

func Propagate(activeRepo *config.ActiveRepo, tree string, defaultParent string, dryRun bool) error {
	home := activeRepo.Repo.Home
	wood.Debugf("arg %s: %s", color.It(color.Blue, "home"), color.It(color.Green, home))
	wood.Debugf("arg %s: %s", color.It(color.Blue, "tree"), color.It(color.Green, tree))
	wood.Debugf("arg %s: %s", color.It(color.Blue, "defaultParent"), color.It(color.Green, defaultParent))
	if dryRun {
		wood.Warnf("arg %s: %s", color.It(color.Blue, "dryRun"), color.It(color.Green, "true"))
	} else {
		wood.Debugf("arg %s: %s", color.It(color.Blue, "dryRun"), color.It(color.Green, "false"))
	}

	treeBranches, ok := activeRepo.FindTree(tree)
	if !ok {
		return errors.Errorf("no tree found: %s", color.It(color.Green, tree))
	}

	branches, names := calcBranches(treeBranches, defaultParent)
	wood.Infof("Propagation: %s", strings.Join(names, " => "))

	gx := gitc{DryRun: dryRun, Home: home}

	parent := branches[0].Name
	for _, branch := range branches[1:] {
		wood.Prefix(branch.Name)

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

func calcBranches(treeBranches []config.TreeBranch, start string) ([]config.TreeBranch, []string) {
	branches := make([]config.TreeBranch, 0, len(treeBranches))
	names := make([]string, 0, len(treeBranches))
	pending := start != ""

	for _, item := range treeBranches {
		if strings.Contains(item.Name, start) {
			pending = true
		}
		if !pending {
			continue
		}
		branches = append(branches, item)
		names = append(names, color.It(color.Green, item.Name))
	}

	return branches, names
}

// gitc wrapper around git cli
type gitc struct {
	Home   string
	DryRun bool
}

func (x *gitc) checkout(branch string) error {
	err := x.run("checkout", branch)
	if err != nil {
		return errors.Wrap(err, "checkout failed")
	}
	wood.Infof("Checked out %s", color.It(color.Green, branch))
	return nil
}

func (x *gitc) rebase(branch string, parent string) error {
	wood.Infof("Rebasing %s with %s", color.It(color.Green, branch), color.It(color.Cyan, parent))
	err := x.run("rebase", "-i", parent)
	if err != nil {
		return errors.Wrap(err, "rebase failed")
	}

	wood.Infof("Rebased %s with %s", color.It(color.Green, branch), color.It(color.Cyan, parent))
	return nil
}

func (x *gitc) pullRebase(branch string, parent string) error {
	err := x.run("pull", "--rebase", parent)
	if err != nil {
		return errors.Wrap(err, "pulled --rebase failed")
	}
	wood.Infof("Pulled %s and rebased %s", color.It(color.Green, branch), color.It(color.Cyan, parent))
	return nil
}

// Inject replaces placeholders in the template and writes the result into the output file
func (x *gitc) run(args ...string) error {
	if x.DryRun {
		return nil
	}

	opts := proc.RunOpts{Dir: x.Home}
	_, _, err := proc.Run("git", opts, args...)
	if err != nil {
		return errors.Wrap(err, "git call failed")
	}

	return nil
}
