package cmd

import (
	"fmt"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/cmd/clientfactory"
	"hass/internal/config"
)

func init() {
	rootCmd.AddCommand(sceneCmd)
}

var sceneCmd = &cobra.Command{
	Use:   "scene <name>",
	Short: "Various light arrangements",
	Long:  "Activate home lights based on different scenarios",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clientfactory.NewHassClient(endpoint)
		if err != nil {
			wood.Fatalf("Failed to create HASS API client: %v", err)
		}

		cm, err := config.NewConfigManager()
		if err != nil {
			return err
		}

		sm := cm.Scenes()

		if ok, err := requireSingleArg(args, func() error {
			entities := sm.List()
			for _, entity := range entities {
				fmt.Println(entity)
			}
			fmt.Println()
			return cmd.Help()
		}); ok || err != nil {
			return err
		}

		name := args[0]
		if ok := sm.HasScene(name); ok {
			return sm.Execute(client, name)
		}
		return nil
	},
}
