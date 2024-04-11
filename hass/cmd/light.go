package cmd

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/internal/config"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(lightCmd)
	lightCmd.AddCommand(lightOffCmd)
	lightOffCmd.Flags().StringP("name", "n", "", "Light name; omit for all lights")
}

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Basic light control",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var lightOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn off light",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := hass.New(endpoint)
		if err != nil {
			wood.Fatalf("Failed to create HASS API client: %v", err)
		}

		sm, err := config.NewSceneManager()
		if err != nil {
			return err
		}

		name := must(cmd.Flags().GetString("name"))
		if name == "" {
			// turn off all lights
			lights := sm.ListLights()
			for _, light := range lights {
				lightId := sm.GetLightId(light)
				if err := client.LightOff(lightId); err != nil {
					wood.Warnf("Failed to turn off light: %s", lightId)
				}
			}
			return nil
		}

		entityId := sm.Light(name)
		return client.LightOff(entityId)
	},
}
