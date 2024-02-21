package config

type Config struct {
	Lights  map[string]string  `yaml:"lights"`
	Success map[string][]Scene `yaml:"scenes"`
}

type Scenes struct {
}

type Scene struct {
	Light      string `yaml:"light"`
	State      string `yaml:"state"`
	Color      string `yaml:"color"`
	Flash      string `yaml:"flash"`
	Duration   string `yaml:"duration"`
	Brightness int    `yaml:"brightness"`
}
