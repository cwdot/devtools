package configmanager

import (
	"embed"
	"hass/internal/config"
	"hass/internal/managers/lightmanager"
	"hass/internal/managers/scenemanager"
	"hass/internal/managers/speakmanager"
	"os"
	"path/filepath"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

//go:embed base-scenes.yaml
var e embed.FS

type internalConfig struct {
	Lights map[string]string               `yaml:"lights"`
	Speak  map[string]config.SpeakerTarget `yaml:"speak"`
	Scenes map[string][]EntityConfig       `yaml:"scenes"`
}

func New() (*ConfigManager, error) {
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

	var originalConfig *config.Config
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

func mergeConfigs(config *config.Config, originalConfig *config.Config) *config.Config {
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

func readConfig(contents []byte) (*config.Config, error) {
	var ic internalConfig
	if err := yaml.Unmarshal(contents, &ic); err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	typedScenes := make(map[string][]config.Entity)
	for k, v := range ic.Scenes {
		typedEntities := make([]config.Entity, 0, len(v))
		for _, entity := range v {
			typedEntity, err := entity.GetEntity()
			if err != nil {
				panic(errors.Wrapf(err, "error getting entity: %s", k))
			}
			typedEntities = append(typedEntities, typedEntity)
		}
		typedScenes[k] = typedEntities
	}

	return &config.Config{
		Lights: ic.Lights,
		Scenes: typedScenes,
		Speak:  ic.Speak,
	}, nil
}

type ConfigManager struct {
	config *config.Config
}

func (c *ConfigManager) Scenes() *scenemanager.SceneManager {
	lm := c.Lights()
	return scenemanager.New(c.config.Scenes, lm)
}

func (c *ConfigManager) Lights() *lightmanager.LightManager {
	return lightmanager.New(c.config.Lights)
}

func (c *ConfigManager) Speaker() *speakmanager.SpeakManager {
	return speakmanager.New(c.config.Speak)
}
