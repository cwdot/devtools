package jiraprovider

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"

	"gitter/internal/config"
)

func GetIssues(config *config.JiraConfig, ids ...string) (map[string]*jira.Issue, error) {
	m := make(map[string]*jira.Issue)
	if config == nil {
		return m, nil
	}

	tp := &jira.PATAuthTransport{
		Token: config.Token,
	}
	jiraClient, _ := jira.NewClient(tp.Client(), config.Domain)
	jql := fmt.Sprintf("key in (%s)", strings.Join(ids, ","))
	issues, _, err := jiraClient.Issue.Search(jql, &jira.SearchOptions{
		StartAt:       0,
		MaxResults:    25,
		Expand:        "",
		Fields:        []string{"status"},
		ValidateQuery: "",
	})
	if err != nil {
		fmt.Println("JQL:", jql)
		return nil, err
	}

	for _, issue := range issues {
		m[issue.Key] = &issue
	}
	return m, nil
}
