package cmd

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/internal/config"
	"hass/internal/hass"
)

func init() {
	rootCmd.AddCommand(speakCmd)
	speakCmd.Flags().StringP("group", "g", "", "Group")
	speakCmd.Flags().StringP("message", "m", "", "Message")
}

var speakCmd = &cobra.Command{
	Use:   "speak",
	Short: "TTS",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := hass.New(endpoint)
		if err != nil {
			wood.Fatalf("Failed to create HASS API client: %v", err)
		}

		cm, err := config.NewConfigManager()
		if err != nil {
			return err
		}

		group := must(cmd.Flags().GetString("group"))
		message := must(cmd.Flags().GetString("message"))

		return cm.Speaker().Speak(client, group, message)
	},
}
