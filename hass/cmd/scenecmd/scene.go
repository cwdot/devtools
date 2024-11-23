package scenecmd

import (
	"fmt"
	"github.com/pkg/errors"
	"hass/internal/managers/configmanager"
	"hass/internal/managers/scenemanager"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/cmd/clientfactory"
)

func NewSceneCmd(endpoint string) *cobra.Command {
	return &cobra.Command{
		Use:   "scene <name>",
		Short: "Various light arrangements",
		Long:  "Activate home lights based on different scenarios",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hc, err := clientfactory.NewHassClient(endpoint)
			if err != nil {
				wood.Fatalf("Failed to create HASS API client: %v", err)
			}

			mc, err := clientfactory.NewMqttClient()
			if err != nil {
				wood.Fatalf("create MQTT client: %v", err)
			}
			if err := mc.Connect(); err != nil {
				wood.Fatalf("connect to MQTT broker: %v", err)
			}
			defer mc.Disconnect()

			cm, err := configmanager.New()
			if err != nil {
				return err
			}

			sm := cm.Scenes()
			switch len(args) {
			case 0: // print scenes
				printScenes(sm)
				return nil
			case 1: // one light
				entityId := args[0]
				if !sm.HasScene(entityId) {
					printScenes(sm)
					return errors.Errorf("not found: %v", entityId)
				}

				return sm.Execute(hc, mc, entityId)
			}
			return nil
		},
	}
}

func printScenes(sm *scenemanager.SceneManager) {
	entities := sm.List()
	for _, entity := range entities {
		fmt.Println(entity)
	}
	fmt.Println()
	fmt.Println()
}
