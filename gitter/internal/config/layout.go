package config

import tw "github.com/olekukonko/tablewriter"

type Column struct {
	Width int        `yaml:"width,omitempty"`
	Kind  ColumnKind `yaml:"kind"`
	Title string     `yaml:"title"`
	Wrap  WrapStyle  `yaml:"wrap,omitempty"`
	Style int        `yaml:"style,omitempty"`
	Color int        `yaml:"color,omitempty"`
}

type ColumnKind string

type WrapStyle string

const (
	Active          ColumnKind = "active"
	LastHash        ColumnKind = "hash"
	LastHashShort   ColumnKind = "hash-short"
	Project         ColumnKind = "project"
	Name            ColumnKind = "name"
	Description     ColumnKind = "description"
	LastCommitted   ColumnKind = "last-committed"
	CommittedDate   ColumnKind = "committed-date"
	RelDate         ColumnKind = "rel-date"
	RootDrift       ColumnKind = "root-drift"
	RootDriftDesc   ColumnKind = "root-drift-desc"
	RootTracking    ColumnKind = "root-tracking"
	RemoteDrift     ColumnKind = "remote-drift"
	RemoteDriftDesc ColumnKind = "remote-drift-desc"
	RemoteTracking  ColumnKind = "remote-tracking"
	JiraStatus      ColumnKind = "jira-status"
	Links           ColumnKind = "links"

	Nothing  WrapStyle = ""
	Truncate WrapStyle = "truncate"
)

func DefaultLayout() []Column {
	return []Column{
		{1, Active, "*", "", tw.Bold, tw.FgHiGreenColor},
		{20, Project, "PROJECT", "", tw.Normal, tw.Normal},
		{30, Name, "NAME", "truncate", tw.Normal, tw.Normal},
		{30, Description, "DESCRIPTION", "truncate", tw.Normal, tw.Normal},
		{30, LastCommitted, "LAST COMMITTED", "truncate", tw.Normal, tw.Normal},
		{20, CommittedDate, "COMMITTED DATE", "", tw.Normal, tw.Normal},
		{8, RelDate, "REL DATE", "", tw.Normal, tw.Normal},
		// {12, RootDriftDesc, "R DRIFT", "", tw.Normal, tw.Normal},
		// {14, RootTracking, "R TRACKING", "truncate", tw.Normal, tw.Normal},
		// {12, RemoteDriftDesc, "REM DRIFT", "", tw.Normal, tw.Normal},
		// {14, RemoteTracking, "REM TRACKING", "truncate", tw.Normal, tw.Normal},
		{20, JiraStatus, "Jira Status", "truncate", tw.Normal, tw.Normal},
		{30, Links, "LINKS", "", tw.Normal, tw.Normal},
	}
}
