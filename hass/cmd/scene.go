package cmd

import (
	"github.com/spf13/cobra"
)

var brightness int

func init() {
	rootCmd.AddCommand(sceneCmd)
	rootCmd.Flags().StringP("scene", "s", "", "Light brightness")

	rootCmd.PersistentFlags().IntVarP(&brightness, "brightness", "b", lowBrightness, "Light brightness")
}

//
//var sceneSuccessCmd = &cobra.Command{
//	Use:   "success",
//	Short: "Build success",
//	Long:  "Activate bloom and bedroom lights with green",
//	Run: func(cmd *cobra.Command, args []string) {
//		must(client.LightOn(bloom, hass.Green(), hass.LongFlash(), hass.Brightness(brightness)))
//		must(client.LightOn(bedroom, hass.Green(), hass.ShortFlash(), hass.Brightness(lowBrightness)))
//	},
//}
//
//var sceneFailureCmd = &cobra.Command{
//	Use:   "failure",
//	Short: "Build failure",
//	Long:  "Activate bloom and bedroom lights with red",
//	Run: func(cmd *cobra.Command, args []string) {
//		must(client.LightOn(bloom, hass.Red(), hass.LongFlash(), hass.Brightness(brightness)))
//		must(client.LightOn(bedroom, hass.Red(), hass.ShortFlash(), hass.Brightness(lowBrightness)))
//	},
//}
//
//var sceneResetCmd = &cobra.Command{
//	Use:   "reset",
//	Short: "Turn off alert light",
//	Long:  "Turn off bloom and bedroom lights",
//	Run: func(cmd *cobra.Command, args []string) {
//		must(client.LightOff(bloom))
//		must(client.LightOff(bedroom))
//	},
//}
//
//var sceneDangerCmd = &cobra.Command{
//	Use:   "danger",
//	Short: "Turn on danger scene",
//	Long:  "Activate bloom and bedroom lights with red",
//	Run: func(cmd *cobra.Command, args []string) {
//		must(client.LightOn(bloom, hass.Yellow(), hass.LongFlash(), hass.Brightness(100)))
//		must(client.LightOn(bedroom, hass.Yellow(), hass.LongFlash(), hass.Brightness(100)))
//	},
//}
//
//var sceneExerciseCmd = &cobra.Command{
//	Use:   "exercise",
//	Short: "Run exercise scene",
//	Long:  "",
//	Run: func(cmd *cobra.Command, args []string) {
//		must(client.LightOn(bloom, hass.Blue(), hass.LongFlash(), hass.Brightness(brightness)))
//		must(client.LightOn(bedroom, hass.Blue(), hass.LongFlash(), hass.Brightness(brightness)))
//	},
//}

var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "Various light arrangements",
	Long:  "Activate home lights based on different scenarios",
	Run: func(cmd *cobra.Command, args []string) {

		scene := must(cmd.Flags().GetString("scene"))

	},
}
