package jiraprovider

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func GetIssues(ids ...string) ([]*jira.Issue, error) {
	jiraClient, _ := jira.NewClient(nil, "https://issues.apache.org/jira/")
	issue, _, _ := jiraClient.Issue.Get(ids[0], nil)

	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
	fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
	fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)

	return []*jira.Issue{issue}, nil
}
