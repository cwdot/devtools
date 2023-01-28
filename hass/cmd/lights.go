package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(lightCmd)
	lightCmd.AddCommand(lightOnCmd)
	lightCmd.AddCommand(lightOffCmd)
}

var lightOnCmd = &cobra.Command{
	Use:   "on",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn("light.hue_bloom_1", hass.Red())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var lightOffCmd = &cobra.Command{
	Use:   "off",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOff("light.hue_bloom_1")
		if err != nil {
			log.Fatal(err)
		}
	},
}

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
