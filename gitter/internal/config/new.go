package config

import (
	"os"
	"path/filepath"
	"strings"

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

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}
	for idx, repo := range data.Repos {
		data.Repos[idx].Home = strings.ReplaceAll(repo.Home, "$HOME", home)
	}

	return data, nil
}

type Config struct {
	Repos   []Repo              `yaml:"repos"`
	Layouts map[string][]Column `yaml:"layouts"`

	Location string `yaml:"-"` // file location done by the reader
}
