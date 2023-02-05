package config

type BranchSet map[string][]Branch
type TreeSet map[string][]TreeBranch

type Repo struct {
	Name       string    `yaml:"name"`
	Home       string    `yaml:"home"`
	RootBranch string    `yaml:"root_branch"`
	BaseLinks  BaseLinks `yaml:"base_links"`
	Active     BranchSet `yaml:"active"`
	Archived   BranchSet `yaml:"archived"`
	Scripts    []Script  `yaml:"scripts"`
	Trees      TreeSet   `yaml:"trees"`
}

type BaseLinks struct {
	PrBase   string `yaml:"pr_base"`
	JiraBase string `yaml:"jira_base"`
	RepoBase string `yaml:"repo_base"`
}

type BranchLinks struct {
	Pr   string `yaml:"pr"`
	Jira string `yaml:"jira"`
}

type Branch struct {
	Name         string      `yaml:"name"`
	RemoteBranch string      `yaml:"remote"`
	Description  string      `yaml:"description"`
	Links        BranchLinks `yaml:"links"`
}

type TreeBranch struct {
	Name string `yaml:"name"`
}

type Script struct {
	Name      string   `yaml:"name"`
	Lifecycle string   `yaml:"lifecycle"`
	Command   string   `yaml:"command"`
	Arguments []string `yaml:"arguments"`
}
