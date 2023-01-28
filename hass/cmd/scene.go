package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(sceneCmd)
	sceneCmd.AddCommand(sceneSuccessCmd)
	sceneCmd.AddCommand(sceneFailureCmd)
	sceneCmd.AddCommand(sceneAlertCmd)
	sceneCmd.AddCommand(sceneExerciseCmd)
	sceneCmd.AddCommand(sceneResetCmd)
}

var sceneSuccessCmd = &cobra.Command{
	Use:   "success",
	Short: "Build success",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Green(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneFailureCmd = &cobra.Command{
	Use:   "failure",
	Short: "Build failure",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Red(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Red(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}

		err = client.LightOn(canvas, hass.Green(), hass.LongFlash(), hass.TurnOff(2))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneExerciseCmd = &cobra.Command{
	Use:   "exercise",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Blue(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}

		err = client.LightOn(desklight, hass.Blue(), hass.LongFlash())
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Turn off the bloom light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOff("light.hue_bloom_1")
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
