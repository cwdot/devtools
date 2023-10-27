package jirap

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/cwdot/go-stdlib/wood"
	"github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"

	"gitter/internal/config"
)

func GetIssuesSlice(config *config.JiraConfig, ids ...string) ([]jira.Issue, error) {
	m, err := GetIssues(config, ids...)
	if err != nil {
		return nil, err
	}
	l := make([]jira.Issue, 0, len(m))
	for _, issue := range m {
		l = append(l, issue)
	}
	return l, nil
}

func GetIssues(config *config.JiraConfig, ids ...string) (map[string]jira.Issue, error) {
	m := make(map[string]jira.Issue)
	if config == nil || !config.Valid() {
		return m, nil
	}

	tp := &jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}
	jiraClient, err := jira.NewClient(tp.Client(), config.Domain)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating client")
	}

	keys := mapset.NewSet(ids...).ToSlice()
	keyStr := strings.Join(keys, ",")
	wood.Debugf("JIRAs to query: %s", keyStr)

	jql := fmt.Sprintf("key in (%s)", keyStr)
	issues, res, err := jiraClient.Issue.Search(jql, &jira.SearchOptions{
		StartAt:       0,
		MaxResults:    len(keys),
		Expand:        "",
		Fields:        []string{"status", "summary"},
		ValidateQuery: "",
	})
	if err != nil {
		fmt.Println("responses:", *res.Response)
		fmt.Println("JQL:", jql)
		return nil, err
	}

	for _, issue := range issues {
		m[issue.Key] = issue
	}
	return m, nil
}
