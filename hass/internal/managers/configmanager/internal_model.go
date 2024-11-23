package configmanager

import (
	"github.com/pkg/errors"
	"hass/internal/config"
)

type EntityConfig struct {
	Light      string `yaml:"light"`
	State      string `yaml:"state"`
	Color      string `yaml:"color"`
	Flash      string `yaml:"flash"`
	Duration   string `yaml:"duration"`
	Brightness int    `yaml:"brightness"`

	EntityId string `yaml:"fan"`

	Queue   string `yaml:"mqtt"`
	Payload any    `yaml:"payload"`
}

func (ec *EntityConfig) GetEntity() (config.Entity, error) {
	switch {
	case ec.Light != "":
		return config.LightEntity{
			Light:      ec.Light,
			State:      ec.State,
			Color:      ec.Color,
			Flash:      ec.Flash,
			Duration:   ec.Duration,
			Brightness: ec.Brightness,
		}, nil
	case ec.Queue != "":
		return config.MqttEntity{
			Mqtt:    ec.Queue,
			Payload: ec.Payload,
		}, nil
	default:
		return nil, errors.Errorf("unknown entity type: %s", ec)
	}
}
