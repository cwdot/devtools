package config

type Config struct {
	Lights map[string]string        `yaml:"lights"`
	Scenes map[string][]Light       `yaml:"scenes"`
	Speak  map[string]SpeakerTarget `yaml:"speak"`
}

type Entity struct {
	Light
}

type Light struct {
	Light      string `yaml:"light"`
	State      string `yaml:"state"`
	Color      string `yaml:"color"`
	Flash      string `yaml:"flash"`
	Duration   string `yaml:"duration"`
	Brightness int    `yaml:"brightness"`
}

type SpeakerTarget struct {
	Players []string `yaml:"players"`
}
