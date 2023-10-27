package gitp

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"gitter/internal/config"
)

func Test_sortBranches(t *testing.T) {
	original := clone()
	refs := clone()

	rand.Seed(1674844778090707000)
	rand.Shuffle(len(refs), func(i, j int) { refs[i], refs[j] = refs[j], refs[i] })

	sortBranches("master", refs)

	for idx, b := range refs {
		fmt.Println(idx, "==>", b.IsHead, " / ", b.Project, " / ", b.BranchConf)
	}

	require.Equal(t, original, refs)
}

func clone() []*GitBranchMetadata {
	newRef := func(project string, name string, head bool) *GitBranchMetadata {
		return &GitBranchMetadata{
			BranchConf: config.Branch{Name: name},
			BranchName: name,
			Project:    project,
		}
	}

	return []*GitBranchMetadata{
		newRef("master", "master", false),
		newRef("bear", "bear", false),
		newRef("bear", "brown", false),
		newRef("bear", "polar", false),
		newRef("foo", "branch1", false),
		newRef("foo", "branch2", false),
		newRef("foo", "branch3", false),
		newRef("foo", "branch4", false),
		newRef("foo", "branch5", true),
		newRef("foo", "branch6", false),
		newRef("foo", "branch7", false),
	}
}
