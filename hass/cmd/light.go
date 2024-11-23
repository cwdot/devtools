package cmd

import (
	"github.com/spf13/cobra"

	"hass/cmd/lightcmd"
)

func init() {
	rootCmd.AddCommand(lightCmd)

	on := lightcmd.NewLightOnCmd()
	lightCmd.AddCommand(on)

	off := lightcmd.NewLightOffCmd()
	lightCmd.AddCommand(off)
}

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Basic light control",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
