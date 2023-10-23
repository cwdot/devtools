package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"

	"gitter/internal/config"
	"gitter/internal/list"

	"github.com/cwdot/go-stdlib/wood"
)

var allBranches bool
var showArchived bool
var noTrackers bool
var layoutName string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allBranches, "all", "a", false, "Show all git branches; default branches defined in config")
	listCmd.Flags().BoolVarP(&showArchived, "archived", "", false, "Show archived branches")
	listCmd.Flags().BoolVarP(&noTrackers, "notrack", "", false, "Hide tracking info for performance")
	listCmd.Flags().StringVarP(&layoutName, "layout", "", "default", "Customize layout; default is 'default'")
}

func open() (*config.ActiveRepo, *git.Repository, []config.Column, error) {
	if underTest {
		env := filepath.Join(homeDir, ".env")
		return config.OpenCustom(env, layoutName, showArchived)
	}
	return config.OpenDefault(layoutName, showArchived)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		activeRepo, g, layout, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		printConfLocation()
		fmt.Printf("  Repo: %v\n", activeRepo.Repo.Home)
		fmt.Println()

		opts := config.PrintOpts{
			Layout:      layout,
			AllBranches: allBranches,
			NoTrackers:  noTrackers,
		}
		list.Print(activeRepo, g, opts)
	},
}
