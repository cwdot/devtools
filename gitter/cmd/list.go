package cmd

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"gitter/internal/config"
	"gitter/internal/list"

	"github.com/cwdot/go-stdlib/wood"
)

var allBranches bool
var showArchived bool
var layoutName string

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allBranches, "all", "a", false, "Show all git branches; default branches defined in config")
	listCmd.Flags().BoolVarP(&showArchived, "archived", "", false, "Show archived branches")
	listCmd.Flags().StringVarP(&layoutName, "layout", "", "default", "Customize layout; default is 'default'")
}

func open() (*config.ActiveRepo, *git.Repository, []config.Column, error) {
	if underTest {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, nil, nil, err
		}
		env := filepath.Join(home, ".env")
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

		list.PrintBranches(activeRepo, g, layout, allBranches)
	},
}
