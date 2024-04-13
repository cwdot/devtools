package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"hass/internal/hass"
)

func NewSceneManager() (*ConfigManager, error) {
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
	Config Config
}

func (c *ConfigManager) Light(name string) string {
	return c.Config.Lights[name]
}

func (c *ConfigManager) Scene(name string) (*ConfigManagerScene, bool) {
	if s, ok := c.Config.Scenes[name]; ok {
		return &ConfigManagerScene{Scenes: s, lights: c.Config.Lights}, true
	}
	return nil, false
}

func (c *ConfigManager) ListLights() []string {
	keys := make([]string, 0, 10)
	for k, _ := range c.Config.Lights {
		keys = append(keys, k)
	}
	return keys
}

func (c *ConfigManager) ListScenes() []string {
	keys := make([]string, 0, 10)
	for k, _ := range c.Config.Scenes {
		keys = append(keys, k)
	}
	return keys
}

func (c *ConfigManager) GetLightId(alias string) string {
	if fullId, ok := c.Config.Lights[alias]; ok {
		return fullId
	}
	return alias
}

type ConfigManagerScene struct {
	lights map[string]string
	Scenes []Light
}

func (c *ConfigManagerScene) Execute(client *hass.Client) error {
	for _, entity := range c.Scenes {
		opts, err := createOpts(entity)
		if err != nil {
			return err
		}

		id := entity.Light
		if fullId, ok := c.lights[entity.Light]; ok {
			id = fullId
		}

		if err := client.Execute(id, opts...); err != nil {
			return err
		}
	}
	return nil
}

func createOpts(entity Light) ([]func(opts *hass.LightOnOpts), error) {
	opts := make([]func(opts *hass.LightOnOpts), 0, 5)

	switch entity.Color {
	case "red":
		opts = append(opts, hass.Red())
	case "green":
		opts = append(opts, hass.Green())
	case "blue":
		opts = append(opts, hass.Blue())
	case "white":
		opts = append(opts, hass.White())
	case "yellow":
		opts = append(opts, hass.Yellow())
	case "":
	default:
		return nil, errors.Errorf("unknown color: %s", entity.Color)
	}

	if entity.Brightness >= 0 {
		opts = append(opts, hass.Brightness(entity.Brightness))
	}

	switch entity.Flash {
	case "long":
		opts = append(opts, hass.LongFlash())
	case "short":
		opts = append(opts, hass.ShortFlash())
	case "":
	default:
		wood.Debugf("unknown flash: %s", entity.Flash)
	}

	return opts, nil
}
