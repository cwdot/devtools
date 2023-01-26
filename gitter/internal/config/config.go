package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func DefaultConfigFile() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}
	configLocation := filepath.Join(home, ".repo_v2.yaml")
	return ReadConfigFile(configLocation)
}

func ReadConfigFile(full string) (*Config, error) {
	contents, err := os.ReadFile(full)
	if err != nil {
		return nil, errors.Wrap(err, "error reading template file")
	}

	data := &Config{}
	err = yaml.Unmarshal(contents, data)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling template yaml")
	}
	data.Location = full
	return data, nil
}

type Config struct {
	Repos []Repo `yaml:"repos"`

	Location string `yaml:"-"` // file location done by the reader
}

type BranchSet map[string][]Branch
type Repo struct {
	Name       string    `yaml:"name"`
	Home       string    `yaml:"home"`
	RootBranch string    `yaml:"root_branch"`
	BaseLinks  BaseLinks `yaml:"base_links"`
	Active     BranchSet `yaml:"active"`
	Archived   BranchSet `yaml:"archived"`
}

type BaseLinks struct {
	PrBase   string `yaml:"pr_base"`
	JiraBase string `yaml:"jira_base"`
	RepoBase string `yaml:"repo_base"`
}

type Links struct {
	Pr   string `yaml:"pr"`
	Jira string `yaml:"jira"`
}

type Branch struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Links       Links  `yaml:"links"`
}
