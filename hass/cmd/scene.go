package cmd

import (
	"github.com/spf13/cobra"
	"hass/internal/hass"
)

var brightness int

func init() {
	rootCmd.AddCommand(sceneCmd)
	sceneCmd.AddCommand(sceneSuccessCmd)
	sceneCmd.AddCommand(sceneFailureCmd)
	sceneCmd.AddCommand(sceneDangerCmd)
	sceneCmd.AddCommand(sceneExerciseCmd)
	sceneCmd.AddCommand(sceneResetCmd)
	sceneCmd.AddCommand(sceneRingCmd)

	rootCmd.PersistentFlags().IntVarP(&brightness, "brightness", "b", lowBrightness, "Light brightness")
}

var sceneSuccessCmd = &cobra.Command{
	Use:   "success",
	Short: "Build success",
	Long:  "Activate bloom light with green",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(bloom, hass.Green(), hass.LongFlash(), hass.Brightness(brightness)))
		must(client.LightOn(bedroom, hass.Green(), hass.ShortFlash(), hass.Brightness(lowBrightness)))
	},
}

var sceneFailureCmd = &cobra.Command{
	Use:   "failure",
	Short: "Build failure",
	Long:  "Activate bloom light with red",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(bloom, hass.Red(), hass.LongFlash(), hass.Brightness(brightness)))
		must(client.LightOn(bedroom, hass.Red(), hass.ShortFlash(), hass.Brightness(lowBrightness)))
	},
}

var sceneDangerCmd = &cobra.Command{
	Use:   "danger",
	Short: "Turn on danger scene",
	Long:  "Activate bloom and canvas lights with red",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(bloom, hass.Yellow(), hass.LongFlash(), hass.Brightness(100)))

		must(client.LightOn(canvas, hass.Yellow(), hass.TurnOff(10), hass.Brightness(100)))
	},
}

var sceneExerciseCmd = &cobra.Command{
	Use:   "exercise",
	Short: "Run exercise scene",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(bloom, hass.Blue(), hass.LongFlash(), hass.Brightness(brightness)))

		must(client.LightOn(desklight, hass.Blue(), hass.LongFlash(), hass.Brightness(brightness)))
	},
}

var sceneRingCmd = &cobra.Command{
	Use:   "ring",
	Short: "Turn on ring light",
	Long:  "Activate ring light with white",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOn(elgato, hass.White(), hass.Brightness(brightness)))
	},
}

var sceneResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Turn off the bloom light",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		must(client.LightOff("light.hue_bloom_1"))
	},
}

var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "Various light arrangements",
	Long:  "Activate home lights based on different scenarios",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
