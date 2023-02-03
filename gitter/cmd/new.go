package cmd

import (
	"github.com/spf13/cobra"
	"gitter/internal/newconf"

	"github.com/cwdot/go-stdlib/wood"
)

func init() {
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		_, g, _, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		err = newconf.Do(g)
		if err != nil {
			wood.Fatal(err)
		}
	},
}
