package jirap

import (
	"regexp"

	"gitter/internal/config"
)

func SafeExtract(conf *config.JiraConfig, branchName string) string {
	if conf == nil {
		return ""
	}
	return Extract(conf.Extraction, branchName)
}

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
