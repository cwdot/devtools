package bazel

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Parse(value string) (*BazelTarget, error) {
	var pkg, target, file string

	// is this a real file?
	if _, err := os.Stat(value); err == nil {
		d, f := filepath.Split(value)
		pkg = filepath.Clean(d)
		file = f
	} else {
		tokens := strings.Split(value, ":")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid format: %s", value)
		}

		if strings.HasSuffix(tokens[1], ".go") {
			file = tokens[1]
		} else {
			target = tokens[1]
		}
	}

	if !strings.HasPrefix(pkg, "//") {
		pkg = "//" + pkg
	}

	return &BazelTarget{
		Package: pkg,
		Target:  target,
		File:    file,
	}, nil
}

type BazelTarget struct {
	Package string
	Target  string
	File    string
}

func (t BazelTarget) String() string {
	if t.File != "" {
		return fmt.Sprintf("%s:%s", t.Package, t.File)
	}
	return fmt.Sprintf("%s:%s", t.Package, t.Target)
}

func Query(query string) ([]*BazelTarget, error) {
	targets := make([]*BazelTarget, 0, 10)

	// Run Bazel query command
	cmd := exec.Command("bazel", "query", query)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running Bazel query:", err)
		return nil, err
	}

	// Parse and process the output
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Process each line of the Bazel query output
		bt, err := Parse(line)
		if err != nil {
			return nil, err
		}

		targets = append(targets, bt)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading Bazel query output:", err)
		return nil, err
	}

	return targets, nil
}
