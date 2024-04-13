package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func NewConfigManager() (*ConfigManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}

	// aka ~/.config/hass/credentials.env
	scenesPath := filepath.Join(home, ".config", "hass", "scenes.yaml")
	f, err := os.Open(scenesPath)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &ConfigManager{config}, nil
}

type ConfigManager struct {
	config Config
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
