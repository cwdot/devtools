package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/cwdot/stdlib-go/wood"
)

//var domain, service string
//var entityId string
//var alias string

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(noderedCmd)

	serviceCmd.Flags().String("domain", "", "Domain like light or button")
	serviceCmd.Flags().String("service", "", "Action to perform like press, turn_on, or turn_off")
	serviceCmd.Flags().String("entity", "", "Home assistant entity id")
	serviceCmd.Flags().String("message", "", "Home assistant message")
	serviceCmd.Flags().StringSlice("k", []string{}, "Custom key/value pairs; separated by a space")

	noderedCmd.Flags().String("alias", "working", "Nodered alias [sleeping, working]")
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

		domain, _ := cmd.Flags().GetString("domain")
		service, _ := cmd.Flags().GetString("service")
		entityId, _ := cmd.Flags().GetString("entity")
		message, _ := cmd.Flags().GetString("message")

		wood.Debugf("Invoked %s with: %s", service, entityId)

		a, _ := cmd.Flags().GetStringSlice("k")
		log.Println(a)

		arguments := map[string]any{
			"message": message,
		}
		err = client.Service(domain, service, arguments)
		if err != nil {
			log.Fatal("errrr")
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

		alias, _ := cmd.Flags().GetString("alias")

		var entityId string
		switch alias {
		case "sleeping":
			entityId = "button.nodered_55f70f069c8768eb"
		case "working":
			entityId = "button.nodered_77db803615a3b240"
		default:
			wood.Fatalf("unknown alias: %v", alias)
		}

		wood.Infof("Calling button.press service (%s) for %s", entityId, alias)

		err = client.ServiceSimple("button", "press", entityId)
		if err != nil {
			log.Fatal(err)
		}
	},
}
