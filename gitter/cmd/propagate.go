package cmd

import (
	"github.com/spf13/cobra"
	"gitter/internal/propagate"

	"github.com/cwdot/go-stdlib/wood"
)

var tree string
var start string
var dryRun bool

func init() {
	rootCmd.AddCommand(propagateCmd)
	propagateCmd.Flags().StringVarP(&tree, "tree", "t", "", "Tree name")
	propagateCmd.Flags().StringVarP(&start, "start", "s", "", "Parent branch; this is the first branch in the chain. Default is master")
	propagateCmd.Flags().BoolVarP(&dryRun, "dryrun", "", false, "Dry run")
}

var propagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if tree == "" {
			panic("Missing --project argument; part of the trees block")
		}

		activeRepo, _, _, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		err = propagate.Propagate(activeRepo, tree, start, dryRun)
		if err != nil {
			wood.Fatal(err)
		}
	},
}
