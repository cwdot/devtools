package cmd

import (
	"fmt"

	"github.com/cwdot/go-stdlib/color"
	"github.com/cwdot/go-stdlib/wood"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"gitter/internal/config"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPrintCmd)
	configCmd.AddCommand(configLayoutCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management; prints location with no arguments",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		printConfLocation()
	},
}

var configPrintCmd = &cobra.Command{
	Use:   "print",
	Short: "Print current repo configuration",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.DefaultConfigFile()
		if err != nil {
			wood.Fatal(err)
		}
		jsonStr, err := yaml.Marshal(conf)
		fmt.Println(string(jsonStr))
	},
}

var configLayoutCmd = &cobra.Command{
	Use:   "layout",
	Short: "Print layout",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.DefaultConfigFile()
		if err != nil {
			wood.Fatal(err)
		}

		for name, layout := range conf.Layouts {
			fmt.Println(color.Red.It(fmt.Sprintf("Layout: %s", name)))
			for _, column := range layout {
				fmt.Printf(" - %s\n", column.Kind)
			}
		}

		fmt.Println()
		fmt.Println(color.Cyan.It("Default layout"))
		for _, column := range config.DefaultLayout() {
			fmt.Printf(" - %s\n", column.Kind)
		}
	},
}

func printConfLocation() {
	conf, err := config.DefaultConfigFile()
	if err != nil {
		wood.Fatal(err)
	}
	fmt.Printf("Config: %v\n", conf.Location)
}
