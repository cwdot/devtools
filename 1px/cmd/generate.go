package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"1px/internal/config"
	"1px/internal/generator"
	"1px/internal/opbridge"
	"github.com/spf13/cobra"

	"github.com/cwdot/go-stdlib/wood"
)

var machine string
var conf string
var output string

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&machine, "machine", "", "", "env machine")
	generateCmd.Flags().StringVarP(&conf, "config", "", "", "Config listing requested credentials")
	generateCmd.Flags().StringVarP(&output, "output", "", defaultCreds(), "Path to new credentials.env file")
}

func defaultCreds() string {
	home, err := os.UserHomeDir()
	if err != nil {
		wood.Fatalf("error finding home dir: %s", err)
	}
	return filepath.Join(home, ".credentials.env")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate credentials.env file",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if machine == "" {
			cmd.Help()
			fmt.Println()
			wood.Fatal("Missing --machine argument")
		}
		if conf == "" {
			cmd.Help()
			fmt.Println()
			wood.Fatal("Missing --conf argument")
		}

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
				wood.Fatal(err)
			}

			for _, entry := range entries {
				vault := entry.Vault.ID
				id := entry.ID

				// only support categories
				if entry.Category != "API_CREDENTIAL" {
					continue
				}

				pairs = append(pairs, generator.Entry{
					Key:     credential.Key,
					Value:   fmt.Sprintf("op://%s/%s/%s", vault, id, credential.Field),
					Comment: entry.Title,
				})
			}
		}

		if len(pairs) == 0 {
			wood.Fatalf("Found zero credentials to export")
		}
		wood.Infof("Found %d credentials to export", len(pairs))

		// write to credentials file
		err = generator.Write(pairs, output)
		if err != nil {
			wood.Fatal(err)
		}
	},
}
