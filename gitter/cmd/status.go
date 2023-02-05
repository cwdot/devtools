package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/color"
	"github.com/spf13/cobra"

	"github.com/cwdot/go-stdlib/proc"
	"github.com/cwdot/go-stdlib/wood"
)

var lifecycle string

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringVarP(&lifecycle, "lifecycle", "l", "status", "Lifecycle to run; default is 'status'")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		activeRepo, _, _, err := open()
		if err != nil {
			wood.Fatal(err)
		}

		if activeRepo.Repo.Scripts == nil {
			wood.Infof("No scripts to run")
			return
		}

		ranScript := false
		for _, script := range activeRepo.Repo.Scripts {
			if script.Lifecycle != lifecycle {
				continue
			}

			name := script.Name
			ranScript = true

			opts := proc.RunOpts{}
			stdout, stderr, err := proc.Run(script.Command, opts, script.Arguments...)
			if err != nil {
				wood.Warn(err)
				if stderr != "" {
					wood.Warn(stderr)
				}
				wood.Errorf("%s failed: %s", name, script.Command)
				break
			}

			icon := fmt.Sprintf("%s%s%s", color.Green, "OK", color.Reset)
			wood.Infof("%s %s: %s", icon, name, stdout)
		}

		if !ranScript {
			wood.Warnf("Found no scripts matching lifecycle: %s", lifecycle)
		}
	},
}
