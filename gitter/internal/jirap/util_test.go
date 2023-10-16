package jirap

import "testing"

func TestExtract(t *testing.T) {
	expr := "([A-Za-z0-9]+-[0-9]+)"
	tests := []struct {
		name       string
		branchName string
		want       string
	}{
		{"blank", "", ""},
		{"jira only", "FOOBAR-123", "FOOBAR-123"},
		{"prefix with jira", "prefix/FOOBAR-123", "FOOBAR-123"},
		{"multiple prefixes with jira", "a/b/c/d/e/f/FOOBAR-123", "FOOBAR-123"},
		{"unknown format 1", "nope1", ""},
		{"unknown format 2", "x-unknown", ""},
		{"master", "master", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Extract(expr, tt.branchName); got != tt.want {
				t.Errorf("findJira() = %v, want %v", got, tt.want)
			}
		})
	}
}
