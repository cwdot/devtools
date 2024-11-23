package scenecmd

import (
	"fmt"
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"hass/internal/managers/configmanager"
	"hass/internal/managers/scenemanager"
	"log"
	"strings"

	"hass/cmd/clientfactory"
)

func NewSceneCmd(endpoint string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scene <name> [action]",
		Short: "Various light arrangements",
		Long:  "Activate home lights based on different scenarios",
		Args:  cobra.MaximumNArgs(1),
		// Args:  cobra.RangeArgs(1, 2),
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

				argSlice, err := cmd.Flags().GetStringSlice("arg")
				if err != nil {
					log.Fatal(err)
				}
				arguments := sliceToMap(argSlice)
				return sm.Execute(hc, mc, entityId, arguments)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceP("arg", "a", []string{}, "Custom key/value argument pairs; separated by a space")

	return cmd
}

func sliceToMap(slice []string) map[string]string {
	m := make(map[string]string)
	for _, s := range slice {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			log.Fatalf("invalid argument: %v", s)
		}
		m[kv[0]] = kv[1]
	}
	return m
}

func printScenes(sm *scenemanager.SceneManager) {
	entities := sm.List()
	for _, entity := range entities {
		fmt.Println(entity)
	}
	fmt.Println()
	fmt.Println()
}
