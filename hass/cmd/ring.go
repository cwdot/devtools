package cmd

import (
	"github.com/spf13/cobra"

	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(ringCmd)
	ringCmd.AddCommand(ringOnCmd)
	ringCmd.AddCommand(ringOffCmd)
}

var ringOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Turn on Elgato",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(elgato, hass.White(), hass.Brightness(brightness)))
	},
}

var ringOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off Elgato",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOff(elgato))
	},
}

var ringCmd = &cobra.Command{
	Use:   "ring",
	Short: "Ring light",
	Long:  "Activate ring light",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
