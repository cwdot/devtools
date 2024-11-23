package scenemanager

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"hass/internal/config"
	"hass/internal/hassclient"
	"hass/internal/managers/lightmanager"
	"hass/internal/managers/mqttmanager"
	"hass/internal/mqttclient"
	"sort"
)

func New(scenes map[string][]config.Entity, lm *lightmanager.LightManager, mm *mqttmanager.MqttManager) *SceneManager {
	return &SceneManager{
		scenes: scenes,
		lm:     lm,
		mm:     mm,
	}
}

type SceneManager struct {
	scenes map[string][]config.Entity
	lm     *lightmanager.LightManager
	mm     *mqttmanager.MqttManager
}

func (c *SceneManager) List() []string {
	keys := make([]string, 0, len(c.scenes))
	for k, _ := range c.scenes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (c *SceneManager) HasScene(name string) bool {
	_, ok := c.scenes[name]
	return ok
}

func (c *SceneManager) Execute(hc *hassclient.Client, mc *mqttclient.Client, entityId string, arguments map[string]string) error {
	invocations, ok := c.scenes[entityId]
	if !ok {
		return errors.Errorf("not found: %v", entityId)
	}

	for _, invocation := range invocations {
		switch t := invocation.(type) {
		case config.LightEntity:
			if err := c.lm.Execute(hc, t); err != nil {
				return err
			}
		case config.MqttEntity:
			if err := c.mm.Execute(mc, t, arguments); err != nil {
				return err
			}
		//case FanEntity:
		//	if err := c.lm.Execute(client, t); err != nil {
		//		return err
		//	}
		default:
			wood.Fatalf("Unknown invocation type: %v", t)
		}
	}
	return nil
}
