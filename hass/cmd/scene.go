package cmd

import (
	scencmd "hass/cmd/scenecmd"
)

func init() {
	rootCmd.AddCommand(scencmd.NewSceneCmd(endpoint))
}
