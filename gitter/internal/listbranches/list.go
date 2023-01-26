package listbranches

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitter/internal/config"
)

type GitBranchMetadata struct {
	Branch         *config.Branch
	Project        string
	IsHead         bool
	Hash           string
	LastCommit     *object.Commit
	TrackingBranch string
}

func SortBranches(layout *config.Layout, g *git.Repository) ([]*GitBranchMetadata, error) {
	iter, err := g.Branches()
	if err != nil {
		log.Panic(err)
	}

	head, err := g.Head()
	if err != nil {
		log.Panic(err)
	}

	refs := make([]*GitBranchMetadata, 0, 20)

	err = iter.ForEach(func(r *plumbing.Reference) error {
		branch, project, _ := layout.FindBranch(r.Name())

		var tracking string
		revision := r.Name()

		lastCommit, err := g.CommitObject(r.Hash())
		if err != nil {
			return err
		}

		gBranch, err := g.Branch(revision.Short())
		if err == nil {
			tracking = fmt.Sprintf("%s/%s", gBranch.Remote, gBranch.Name)
		}

		refs = append(refs, &GitBranchMetadata{
			Branch:         branch,
			Project:        project,
			LastCommit:     lastCommit,
			TrackingBranch: tracking,
			Hash:           r.Hash().String(),
			IsHead:         r.Hash() == head.Hash(),
		})
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Branch == nil && refs[j].Branch == nil {
			return refs[i].Hash < refs[j].Hash
		} else if refs[i].Branch != nil && refs[j].Branch == nil {
			return true
		} else if refs[i].Branch == nil && refs[j].Branch != nil {
			return false
		}

		if refs[i].Project < refs[j].Project {
			return true
		}
		if refs[i].Branch.Name < refs[j].Branch.Name {
			return true
		}
		return false
	})

	return refs, nil
}

func GenerateLinks(base config.BaseLinks, links config.Links) string {
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
