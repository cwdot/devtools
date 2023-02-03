package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cwdot/go-stdlib/wood"
)

func init() {
	rootCmd.AddCommand(titleCmd)
}

var titleCmd = &cobra.Command{
	Use:   "title",
	Short: "Print title based on JIRA and description of current branch",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		activeRepo, g, _, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		head, err := g.Head()
		if err != nil {
			wood.Fatal(err)
		}

		shortName := head.Name().Short()
		ref, ok := activeRepo.FindBranch(shortName)
		if !ok {
			wood.Fatal("Unknown branch: %v", shortName)
		}

		parts := make([]string, 0, 3)
		branch := ref.Branch

		if branch.Links.Jira != "" {
			parts = append(parts, branch.Links.Jira)
		}
		if branch.Description != "" {
			parts = append(parts, branch.Description)
		}
		fmt.Println(strings.Join(parts, " "))
	},
}
