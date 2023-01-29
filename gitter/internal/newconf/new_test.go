package newconf

import "testing"

func Test_findJira(t *testing.T) {
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
			if got := findJira(tt.branchName); got != tt.want {
				t.Errorf("findJira() = %v, want %v", got, tt.want)
			}
		})
	}
}
