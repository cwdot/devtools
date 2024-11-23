package scenemanager

import (
	"encoding/json"
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"hass/internal/config"
	"hass/internal/hassclient"
	"hass/internal/managers/lightmanager"
	"hass/internal/mqttclient"
	"sort"
)

func New(scenes map[string][]config.Entity, lm *lightmanager.LightManager) *SceneManager {
	return &SceneManager{
		scenes: scenes,
		lm:     lm,
	}
}

type SceneManager struct {
	lm     *lightmanager.LightManager
	scenes map[string][]config.Entity
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

func (c *SceneManager) Execute(hc *hassclient.Client, mc *mqttclient.Client, entityId string) error {
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
			if err := executeMqtt(mc, t); err != nil {
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

func executeMqtt(mc *mqttclient.Client, invocation config.MqttEntity) error {
	wood.Infof("Invoking MQTT: %s", invocation.Mqtt)

	var payload string
	switch t := invocation.Payload.(type) {
	case []any:
		p, err := json.Marshal(t[0])
		if err != nil {
			return errors.Wrap(err, "marshalling payload")
		}
		payload = string(p)
	}

	return mc.Publish(invocation.Mqtt, payload, false)
}
