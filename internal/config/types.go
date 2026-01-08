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

type ModbusTCPConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	UnitID        int    `yaml:"unitId"`
	Timeout       string `yaml:"timeout"`
	PollInterval  string `yaml:"pollInterval"`
	RetryInterval string `yaml:"retryInterval"`
	ByteOrder     string `yaml:"byteOrder"`
	WordOrder     string `yaml:"wordOrder"`
}

type ModbusRTUConfig struct {
	Device       string `yaml:"device"`
	BaudRate     int    `yaml:"baudRate"`
	DataBits     int    `yaml:"dataBits"`
	Parity       string `yaml:"parity"`
	StopBits     int    `yaml:"stopBits"`
	UnitID       int    `yaml:"unitId"`
	Timeout      string `yaml:"timeout"`
	PollInterval string `yaml:"pollInterval"`
}

func (c *ModbusRTUConfig) Validate() error {
	return nil
}

type Variables struct {
	Name     string   `yaml:"name"`
	Source   string   `yaml:"source"`
	DataType string   `yaml:"dataType"`
	Tags     []string `yaml:"tags"`
}
