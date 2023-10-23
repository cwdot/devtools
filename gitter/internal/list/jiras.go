package list

import (
	"github.com/cwdot/go-stdlib/wood"
	"github.com/pkg/errors"

	"gitter/internal/config"
	"gitter/internal/jirap"
	"gitter/internal/providers/gitprovider"
	"gitter/internal/providers/jiraprovider"
)

func getBranchJiras(jiraConfig *config.JiraConfig, rows []*gitprovider.GitBranchMetadata) (map[string]string, error) {
	// branch to status
	m := make(map[string]string)

	if jiraConfig == nil {
		return m, nil
	}

	// jira to branch
	jiraToBranch := make(map[string]string)

	jiras := make([]string, 0, 10)
	for _, row := range rows {
		if row.BranchConf.Jira != "" {
			jiras = append(jiras, row.BranchConf.Jira)
			jiraToBranch[row.BranchConf.Jira] = row.BranchName
		} else if key := jirap.SafeExtract(jiraConfig, row.BranchName); key != "" {
			jiras = append(jiras, key)
			jiraToBranch[key] = row.BranchName
		}
	}

	issues, err := jiraprovider.GetIssues(jiraConfig, jiras...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query jira")
	}

	for _, issue := range issues {
		branch, ok := jiraToBranch[issue.Key]
		if !ok {
			wood.Warnf("No branch found for Jira %s", issue.Key)
			continue
		}
		m[branch] = issue.Fields.Status.Name
	}

	return m, nil
}
