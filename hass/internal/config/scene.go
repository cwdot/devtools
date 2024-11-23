package config

import (
	"github.com/pkg/errors"

	"hass/internal/hassclient"
)

type SceneManager struct {
	lm     *LightManager
	scenes map[string][]Light
}

func (c *SceneManager) List() []string {
	keys := make([]string, 0, len(c.scenes))
	for k, _ := range c.scenes {
		keys = append(keys, k)
	}
	return keys
}

func (c *SceneManager) HasScene(name string) bool {
	_, ok := c.scenes[name]
	return ok
}

func (c *SceneManager) Execute(client *hassclient.Client, entityId string) error {
	lights, ok := c.scenes[entityId]
	if !ok {
		return errors.Errorf("not found: %v", entityId)
	}

	for _, light := range lights {
		opts, err := createLightOpts(light)
		if err != nil {
			return err
		}

		id := light.Light
		if fullId := c.lm.GetLightId(light.Light); fullId != "" {
			id = fullId
		}

		if err := client.LightOn(id, opts...); err != nil {
			return err
		}
	}
	return nil
}
