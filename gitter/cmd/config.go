package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitter/internal/config"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configPrintCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management; prints location with no arguments",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.DefaultConfigFile()
		if err != nil {
			wood.Fatal(err)
		}
		fmt.Println(conf.Location)
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
