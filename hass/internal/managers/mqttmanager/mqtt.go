package mqttmanager

import (
	"encoding/json"
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"hass/internal/config"
	"hass/internal/mqttclient"
	"strings"
)

func New() *MqttManager {
	return &MqttManager{}
}

type MqttManager struct {
}

func (c *MqttManager) Execute(mc *mqttclient.Client, invocation config.MqttEntity, arguments map[string]string) error {
	wood.Infof("Invoking MQTT: %s", invocation.Mqtt)

	var payload string
	switch t := invocation.Payload.(type) {
	case []any:
		p, err := json.Marshal(t[0])
		if err != nil {
			return errors.Wrap(err, "marshalling payload")
		}
		payload = string(p)
		payload = replaceArguments(payload, arguments)
		wood.Debugf("Invoking MQTT: %s", payload)
	}

	return mc.Publish(invocation.Mqtt, payload, false)
}

func replaceArguments(payload string, arguments map[string]string) string {
	idx := 0

	replacements := make(map[string]string)

	for {
		start := strings.Index(payload[idx:], "${")
		if start == -1 {
			break
		}
		end := strings.Index(payload[start:], "}")
		if end == -1 {
			break
		}
		idx = end + start + 1

		original := payload[start:idx]
		text := original[2 : len(original)-1]
		tokens := strings.SplitN(text, ":", 2)

		key := tokens[0]
		value, ok := arguments[key]
		if ok {
			replacements[original] = value
			continue
		}

		if len(tokens) == 2 {
			// use default value
			replacements[original] = tokens[1]
		}
	}

	for k, v := range replacements {
		payload = strings.ReplaceAll(payload, k, v)
	}

	return payload
}
