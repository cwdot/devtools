package cmd

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cwdot/go-stdlib/wood"
)

var domain, service, entityId string
var alias string

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(noderedCmd)

	serviceCmd.Flags().StringVar(&domain, "domain", "", "Domain like light or button")
	serviceCmd.Flags().StringVar(&service, "service", "", "Action to perform like press, turn_on, or turn_off")
	serviceCmd.Flags().StringVar(&entityId, "entity", "", "Home assistant entity id")
	noderedCmd.Flags().StringVar(&alias, "alias", "working", "Nodered alias [sleeping, working]")
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Build success",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Flags().Parse(args)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			wood.SetLevel(logrus.DebugLevel)
		}

		wood.Debugf("Invoked %s with: %s", service, entityId)
		err = client.ServiceSimple(domain, service, entityId)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var noderedCmd = &cobra.Command{
	Use:   "nodered",
	Short: "Special nodered buttons",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Flags().Parse(args)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			wood.SetLevel(logrus.DebugLevel)
		}

		var entityId string
		switch alias {
		case "sleeping":
			entityId = "button.nodered_55f70f069c8768eb"
		case "working":
			entityId = "button.nodered_77db803615a3b240"
		default:
			panic(fmt.Sprintf("unknown alias: %v", alias))
		}

		wood.Infof("Calling button.press service (%s) for %s", entityId, alias)

		err = client.ServiceSimple("button", "press", entityId)
		if err != nil {
			log.Fatal(err)
		}
	},
}
