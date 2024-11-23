package cmd

import (
	"fmt"
	"os"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/spf13/cobra"

	"hass/cmd/clientfactory"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "hass",
	Short: "Hass",
	Long:  `Home assistant tool`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			wood.SetLevel(wood.DebugLevel)
		}
	},
}

// var client *hassclient.Client
var endpoint string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "HASS endpoint")
}

func Execute() {
	preloadClients()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func preloadClients() {
	_, err := clientfactory.NewHassClient(endpoint)
	if err != nil {
		wood.Fatalf("Failed to create HASS API client: %v", err)
	}

	if _, err := clientfactory.NewMqttClient(); err != nil {
		wood.Fatalf("Failed to create MQTT API client: %v", err)
	}
}
