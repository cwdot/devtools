package lightcmd

import (
	"fmt"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/cmd/clientfactory"
	"hass/internal/config"
)

func NewLightOnCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "on [name]",
		Short: "Turn on light",
		Long:  "Turn on light. If no name is provided, all lights will be turned on.",
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
			fmt.Println("light_on.go")
			lm := cm.Lights()
			switch len(args) {
			case 0: // all lights
				lights := lm.List()
				for _, light := range lights {
					lightId := lm.GetLightId(light)
					if err := client.LightOn(lightId); err != nil {
						wood.Warnf("Failed to turn on light: %s => %v", lightId, err)
					}
				}
				wood.Infof("Turned on all lights (%d)", len(lights))
				return nil
			case 1: // one light
				entityId := lm.GetLightId(args[0])
				return client.LightOn(entityId)
			}
			return nil
		},
	}
}
