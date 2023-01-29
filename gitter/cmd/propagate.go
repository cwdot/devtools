package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitter/internal/propagate"
)

var propagateProject string
var dryRun bool

func init() {
	rootCmd.AddCommand(propagateCmd)
	propagateCmd.Flags().StringVarP(&propagateProject, "project", "", "", "Customize layout; default is 'default'")
	propagateCmd.Flags().BoolVarP(&dryRun, "dryrun", "", false, "Dry run")
}

var propagateCmd = &cobra.Command{
	Use:   "propagate",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if propagateProject == "" {
			panic("Missing --project argument; part of the trees block")
		}

		activeRepo, _, _, err := open()
		if err != nil {
			log.Panic(err)
		}

		err = propagate.Propagate(activeRepo, "propagate", dryRun)
		if err != nil {
			log.Panic(err)
		}
	},
}
