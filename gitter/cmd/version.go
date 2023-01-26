package cmd

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gitter",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var commit string
		var ts time.Time
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					commit = setting.Value
				case "vcs.time":
					ts, _ = time.Parse(time.RFC3339, setting.Value)
				}
			}
		}
		if commit != "" {
			fmt.Println("commit:   ", commit)
		}
		if !ts.IsZero() {
			fmt.Println("timestamp:", ts.Local().String())
		}
	},
}
