package clientfactory

import (
	"hass/internal/mqttclient"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func NewMqttClient() (*mqttclient.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}

	// aka ~/.config/hass/credentials.env
	credentialsPath := filepath.Join(home, ".config", "hass", "credentials.env")
	env, err := LoadAndValidateEnv(credentialsPath, []string{"MQTT_USERNAME", "MQTT_PASSWORD", "MQTT_BROKER", "MQTT_CLIENT_ID"})
	if err != nil {
		return nil, errors.Wrapf(err, "env validation")
	}

	return mqttclient.New(mqttclient.Config{
		Broker:   env["MQTT_BROKER"],
		ClientID: env["MQTT_CLIENT_ID"],
		Username: env["MQTT_USERNAME"],
		Password: env["MQTT_PASSWORD"],
	})
}
