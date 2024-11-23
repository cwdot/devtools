package config

import "fmt"

type Config struct {
	Lights map[string]string
	Scenes map[string][]Entity
	Speak  map[string]SpeakerTarget
}

func (c *Config) Summary() string {
	return fmt.Sprintf("Lights: %v, Scenes: %v, Speak: %v", len(c.Lights), len(c.Scenes), len(c.Speak))
}

type Entity interface {
}

type LightEntity struct {
	Entity
	Light      string `yaml:"light"`
	State      string `yaml:"state"`
	Color      string `yaml:"color"`
	Flash      string `yaml:"flash"`
	Duration   string `yaml:"duration"`
	Brightness int    `yaml:"brightness"`
}

type MqttEntity struct {
	Entity
	Mqtt           string
	Payload        any
	ValidArguments map[string][]string
}

type FanEntity struct {
	Entity
	EntityId string `yaml:"fan"`
}

type SpeakerTarget struct {
	Players []string `yaml:"players"`
}
