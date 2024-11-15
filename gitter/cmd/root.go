package cmd

import (
	"fmt"
	"os"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"
)

var verbose bool
var underTest bool
var quickPrint bool

var rootCmd = &cobra.Command{
	Use:   "gitter",
	Short: "Git repository management",
	Long:  `Manage multiple git branches`,
	Run: func(cmd *cobra.Command, args []string) {
		if quickPrint {
			printConfLocation()
		} else {
			listCmd.Run(cmd, args)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			wood.SetLevel(wood.DebugLevel)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&underTest, "test", "t", false, "Test with env repo")
	rootCmd.Flags().BoolVarP(&quickPrint, "", "c", false, "Print config location")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
