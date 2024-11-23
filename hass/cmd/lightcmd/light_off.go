package lightcmd

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/cmd/clientfactory"
	"hass/internal/config"
)

func NewLightOffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "off [name]",
		Short: "Turn off light",
		Long:  "Turn off light. If no name is provided, all lights will be turned off.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientfactory.NewHassClient("")
			if err != nil {
				wood.Fatalf("Failed to create HASS API client: %v", err)
			}

			cm, err := config.NewConfigManager()
			if err != nil {
				return err
			}

			lm := cm.Lights()
			switch len(args) {
			case 0: // all lights
				lights := lm.List()
				for _, light := range lights {
					lightId := lm.GetLightId(light)
					if err := client.LightOff(lightId); err != nil {
						wood.Warnf("Failed to turn off light: %s => %v", lightId, err)
					}
				}
				wood.Infof("Turned off all lights (%d)", len(lights))
				return nil
			case 1: // one light
				entityId := lm.GetLightId(args[0])
				return client.LightOff(entityId)
			}
			return nil
		},
	}
}
