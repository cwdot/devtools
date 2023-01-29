package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"hass/internal/hass"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "hass",
	Short: "Hass",
	Long:  `Home assistant tool`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var client *hass.Client

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
}

func Execute() {
	var err error
	client, err = hass.New()
	if err != nil {
		log.Fatalf("Failed to create HASS API client: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
