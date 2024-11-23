package config

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

//go:embed base-scenes.yaml
var e embed.FS

func NewConfigManager() (*ConfigManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}

	b, err := e.ReadFile("base-scenes.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "reading embedded")
	}

	baseConfig, err := readConfig(b)
	if err != nil {
		return nil, errors.Wrap(err, "reading embedded")
	}

	var originalConfig *Config
	scenesPath := filepath.Join(home, ".config", "hass", "scenes.yaml")
	if _, err := os.Stat(scenesPath); err == nil {
		b, err := os.ReadFile(scenesPath)
		if err != nil {
			wood.Infof("Scene config path: %v", scenesPath)
			return nil, err
		}

		originalConfig, err = readConfig(b)
		if err != nil {
			return nil, errors.Wrap(err, "reading file")
		}
	}

	config := mergeConfigs(baseConfig, originalConfig)

	return &ConfigManager{config}, nil
}

func mergeConfigs(config *Config, originalConfig *Config) *Config {
	for k, v := range originalConfig.Scenes {
		config.Scenes[k] = v
	}
	for k, v := range originalConfig.Lights {
		config.Lights[k] = v
	}
	for k, v := range originalConfig.Speak {
		config.Speak[k] = v
	}
	return config
}

func readConfig(contents []byte) (*Config, error) {
	var config Config
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}
	return &config, nil
}

type ConfigManager struct {
	config *Config
}

func (c *ConfigManager) Scenes() *SceneManager {
	lm := c.Lights()
	return &SceneManager{scenes: c.config.Scenes, lm: lm}
}

func (c *ConfigManager) Lights() *LightManager {
	return &LightManager{
		aliases: c.config.Lights,
	}
}

func (c *ConfigManager) Speaker() *SpeakManager {
	return &SpeakManager{c.config.Speak}
}
