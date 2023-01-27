package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&underTest, "test", "t", false, "Test with env repo")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
