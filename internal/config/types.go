package config

type Config struct {
	LogLevel  string
	Sources   []Sources
	Variables []Variables
}

type Sources struct {
	Name   string       `yaml:"name"`
	Type   string       `yaml:"type"`
	Config SourceConfig // Interface - will hold concrete type
}

type SourceConfig struct {
	Validate error
}

type Variables struct {
	Name     string   `yaml:"name"`
	Source   string   `yaml:"source"`
	DataType string   `yaml:"dataType"`
	Tags     []string `yaml:"tags"`
}
