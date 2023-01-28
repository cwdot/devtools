package cmd

import (
	"fmt"
	"strings"

	"1px/internal/config"
	"1px/internal/generator"
	"1px/internal/opbridge"
	"github.com/spf13/cobra"
)

var machine string
var conf string
var output string

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&machine, "machine", "", "", "")
	generateCmd.Flags().StringVarP(&conf, "config", "", "", "")
	generateCmd.Flags().StringVarP(&output, "output", "", "", "")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate credentials.env file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// read config.yaml
		config, err := config.ReadConfigFile(conf)
		if err != nil {
			panic(err)
		}

		// query op for all
		pairs := make([]generator.Entry, 0, 10)

		for _, credential := range config.Credentials {
			tag := strings.ReplaceAll(credential.Tags, "$MACHINE", machine)

			entries, err := opbridge.List(tag)
			if err != nil {
				panic(err)
			}

			for _, entry := range entries {
				vault := entry.Vault.ID
				id := entry.ID

				pairs = append(pairs, generator.Entry{
					Key:     credential.Key,
					Value:   fmt.Sprintf("op://%s/%s/%s", vault, id, credential.Field),
					Comment: entry.Title,
				})
			}
		}

		if len(pairs) == 0 {
			panic("Found zero credentials to export")
		}
		fmt.Printf("Found %d credentials to export\n", len(pairs))

		// write to credentials file
		err = generator.Write(pairs, output)
		if err != nil {
			panic(err)
		}
	},
}
