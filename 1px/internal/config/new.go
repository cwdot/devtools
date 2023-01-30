package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func ReadConfigFile(full string) (*Config, error) {
	contents, err := os.ReadFile(full)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}

	data := &Config{}
	err = yaml.Unmarshal(contents, data)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling yaml")
	}
	return data, nil
}

type Config struct {
	Credentials []Credential `yaml:"credentials"`
}

type Credential struct {
	Tags string `yaml:"tags"`
	//Key        string `yaml:"key"`         // explicit value
	//Value      string `yaml:"value"`       // explicit value
	KeyLabel   string `yaml:"key_label"`   // read value from 1p
	ValueLabel string `yaml:"value_label"` // read value from 1p
}
