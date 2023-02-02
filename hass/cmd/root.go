package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"hass/internal/hass"

	"github.com/cwdot/go-stdlib/wood"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "hass",
	Short: "Hass",
	Long:  `Home assistant tool`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			wood.SetLevel(logrus.DebugLevel)
		}
	},
}

var client *hass.Client

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
}

func Execute() {
	disabled := os.Getenv("HASS_DISABLED")
	if disabled != "" {
		wood.Infof("HASS_DISABLED env var set; exiting early")
		return
	}

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
