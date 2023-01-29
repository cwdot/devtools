package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cwdot/go-stdlib/wood"
)

var verbose bool
var underTest bool

var rootCmd = &cobra.Command{
	Use:   "gitter",
	Short: "Git repository management",
	Long:  `Manage multiple git branches`,
	Run: func(cmd *cobra.Command, args []string) {
		listCmd.Run(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			wood.SetLevel(logrus.DebugLevel)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&underTest, "test", "t", false, "Test with env repo")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
