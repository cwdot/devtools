package jiraprovider

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"

	"gitter/internal/config"
)

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
	jql := fmt.Sprintf("key in (%s)", strings.Join(ids, ","))
	issues, res, err := jiraClient.Issue.Search(jql, &jira.SearchOptions{
		StartAt:       0,
		MaxResults:    25,
		Expand:        "",
		Fields:        []string{"status"},
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
