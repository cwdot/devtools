package configmanager

import (
	"fmt"
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

	Queue     string `yaml:"mqtt"`
	Payload   any    `yaml:"payload"`
	Arguments any    `yaml:"arguments"`
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
			Mqtt:           ec.Queue,
			Payload:        ec.Payload,
			ValidArguments: parseMqttArguments(ec.Arguments),
		}, nil
	default:
		return nil, errors.Errorf("unknown entity type: %v", ec)
	}
}

func parseMqttArguments(arguments any) map[string][]string {
	m := make(map[string][]string)
	switch t := arguments.(type) {
	case map[string]interface{}:
		for k, mv := range t {
			slice := make([]string, 0)
			for _, v := range mv.([]interface{}) {
				slice = append(slice, fmt.Sprintf("%v", v))
			}
			m[k] = slice
		}
	}
	return m
}
