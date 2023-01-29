package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(alertCmd)
	alertCmd.AddCommand(alertOnCmd)
	alertCmd.AddCommand(alertOffCmd)
	alertCmd.AddCommand(alertFlashCmd)
}

var alertOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Turns on the bloom (alert) light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn("light.hue_bloom_1", hass.Red())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var alertOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turns off the bloom (alert) light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOff("light.hue_bloom_1")
		if err != nil {
			log.Fatal(err)
		}
	},
}

var alertFlashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Turns on the bloom (alert) light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn("light.hue_bloom_1", hass.Red(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}

		err = client.Deactivate("light.hue_bloom_1", 20*time.Second)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Control the bloom (alert) light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
