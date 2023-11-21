package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cwdot/stdlib-go/wood"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/spf13/cobra"

	"bazres/internal/bazel"
)

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringP("lifecycle", "l", "status", "Lifecycle to run; default is 'status'")
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		knownTargets, err := bazel.Query("kind(go_test, //...)")
		if err != nil {
			return err
		}

		var targets []*bazel.BazelTarget
		if len(args) == 0 {
			var inputReader = cmd.InOrStdin()
			s, err := io.ReadAll(inputReader)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %v", err)
			}
			targets, err = processLines(strings.Split(string(s), "\n")...)
		} else {
			targets, err = processLines(args...)
		}
		if err != nil {
			return err
		}
		collate(knownTargets, targets)

		return nil
	},
}

func collate(testTargets []*bazel.BazelTarget, requestedTargets []*bazel.BazelTarget) {
	m := make(map[string]*bazel.BazelTarget)
	for _, target := range testTargets {
		m[target.Package] = target
	}

	// do we have matching test target?
	for _, target := range requestedTargets {
		if testTarget, ok := m[target.Package]; ok {
			eprint(testTarget.Package + ":" + testTarget.Target)
		}
	}
}

func processLines(args ...string) ([]*bazel.BazelTarget, error) {
	targets := mapset.NewSet[*bazel.BazelTarget]() // where T is some concrete comparable type.
	for _, arg := range args {
		if arg == "" {
			continue
		}

		bt, err := bazel.Parse(arg)
		if err != nil {
			wood.Debugf("failed to parse line: %s => %s", arg, err)
			continue
		}

		bazelFile := getBazelFile(bt.Package)
		if bazelFile == "" {
			wood.Debugf("no bazel file found for: %s", bt.Package)
			continue
		}

		valid, err := hasTestTarget(bazelFile)
		if err != nil {
			wood.Debugf("failed to parse bazel file: %s => %s", bazelFile, err)
			return nil, err
		}
		if !valid {
			wood.Tracef("no test target in bazel file: %s", bazelFile)
			continue
		}

		switch {
		case bt.Target != "":
			targets.Add(bt)
		case bt.File != "":
			wood.Debugf("no cannot read bt file: %s", bt.File)
		default:
			wood.Debugf("no target or file found for: %s", arg)
		}
	}
	return targets.ToSlice(), nil
}

func getBazelFile(p string) string {
	p = strings.TrimPrefix(p, "//")
	for _, candidate := range []string{p + "/BUILD.bazel"} {
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			return candidate
		}
	}
	return ""

}
func hasTestTarget(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	lines, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(lines), "\n") {
		if strings.HasPrefix(line, "go_test(") {
			return true, nil
		}
	}

	return false, nil
}

func eprint(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
