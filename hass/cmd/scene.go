package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"hass/internal/hass"
)

var lightBrightness int

func init() {
	rootCmd.AddCommand(sceneCmd)
	sceneCmd.AddCommand(sceneSuccessCmd)
	sceneCmd.AddCommand(sceneFailureCmd)
	sceneCmd.AddCommand(sceneDangerCmd)
	sceneCmd.AddCommand(sceneExerciseCmd)
	sceneCmd.AddCommand(sceneResetCmd)
	sceneCmd.AddCommand(sceneRingCmd)

	rootCmd.PersistentFlags().IntVarP(&lightBrightness, "brightness", "b", 20, "Light brightness")
}

var sceneSuccessCmd = &cobra.Command{
	Use:   "success",
	Short: "Build success",
	Long:  "Activate bloom light with green",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Green(), hass.LongFlash(), hass.Brightness(lightBrightness))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneFailureCmd = &cobra.Command{
	Use:   "failure",
	Short: "Build failure",
	Long:  "Activate bloom light with red",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Red(), hass.LongFlash(), hass.Brightness(lightBrightness))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneDangerCmd = &cobra.Command{
	Use:   "danger",
	Short: "",
	Long:  "Activate bloom and canvas lights with red",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(bloom, hass.Yellow(), hass.LongFlash(), hass.Brightness(100))
		if err != nil {
			log.Fatal(err)
		}

		err = client.LightOn(canvas, hass.Yellow(), hass.TurnOff(10), hass.Brightness(100))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var sceneRingCmd = &cobra.Command{
	Use:   "ring",
	Short: "",
	Long:  "Activate ring light with white",
	Run: func(cmd *cobra.Command, args []string) {
		err := client.LightOn(elgato, hass.White(), hass.LongFlash(), hass.Brightness(lightBrightness))
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
		err := client.LightOn(bloom, hass.Blue(), hass.LongFlash(), hass.Brightness(lightBrightness))
		if err != nil {
			log.Fatal(err)
		}

		err = client.LightOn(desklight, hass.Blue(), hass.LongFlash(), hass.Brightness(lightBrightness))
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
