package jirap

import (
	"regexp"
)

func Extract(expr string, branchName string) string {
	if expr == "" || branchName == "" {
		return ""
	}

	r, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}

	matches := r.FindStringSubmatch(branchName)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}
