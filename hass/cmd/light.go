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
}

var lightCmd = &cobra.Command{
	Use:   "light",
	Short: "Basic light control",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var lightOffCmd = &cobra.Command{
	Use:   "off <name>",
	Short: "Turn off light",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := hass.New(endpoint)
		if err != nil {
			wood.Fatalf("Failed to create HASS API client: %v", err)
		}

		cm, err := config.NewConfigManager()
		if err != nil {
			return err
		}

		lm := cm.Lights()

		if ok, err := requireSingleArg(args, func() error {
			// turn off all lights
			lights := lm.List()
			for _, light := range lights {
				lightId := lm.GetLightId(light)
				if err := client.LightOff(lightId); err != nil {
					wood.Warnf("Failed to turn off light: %s", lightId)
				}
			}
			return cmd.Help()
		}); ok || err != nil {
			return err
		}

		name := args[0]
		entityId := lm.GetLightId(name)
		return client.LightOff(entityId)
	},
}
