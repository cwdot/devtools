package config

type Config struct {
	Lights map[string]string
	Scenes map[string][]Entity
	Speak  map[string]SpeakerTarget
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
	Mqtt    string `yaml:"mqtt"`
	Payload any    `yaml:"payload"`
}

type FanEntity struct {
	Entity
	EntityId string `yaml:"fan"`
}

type SpeakerTarget struct {
	Players []string `yaml:"players"`
}
