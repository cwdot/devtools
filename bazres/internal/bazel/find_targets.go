package bazel

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func Parse(value string) (*BazelTarget, error) {
	tokens := strings.Split(value, ":")
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid format: %s", value)
	}
	if strings.HasSuffix(tokens[1], ".go") {
		return &BazelTarget{
			Package: tokens[0],
			Target:  "",
			File:    tokens[1],
		}, nil
	}
	return &BazelTarget{
		Package: tokens[0],
		Target:  tokens[1],
		File:    "",
	}, nil
}

type BazelTarget struct {
	Package string
	Target  string
	File    string
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
