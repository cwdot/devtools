package cmd

import (
	"fmt"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/internal/config"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(sceneCmd)
	sceneCmd.Flags().StringP("name", "n", "", "Scene name")
}

var sceneCmd = &cobra.Command{
	Use:   "scene",
	Short: "Various light arrangements",
	Long:  "Activate home lights based on different scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := hass.New(endpoint)
		if err != nil {
			wood.Fatalf("Failed to create HASS API client: %v", err)
		}

		cm, err := config.NewConfigManager()
		if err != nil {
			return err
		}

		sm := cm.Scenes()

		scene := must(cmd.Flags().GetString("name"))
		if scene == "" {
			entities := sm.List()
			for _, entity := range entities {
				fmt.Println(entity)
			}
			return cmd.Help()
		}

		if ok := sm.HasScene(scene); ok {
			return sm.Execute(client, scene)
		}
		return nil
	},
}
