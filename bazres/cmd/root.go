package cmd

import (
	"fmt"
	"os"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "bazres",
	Short: "Bazel resolver",
	Long:  `Resolve bazel targets`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			wood.SetLevel(wood.DebugLevel)
		}
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
