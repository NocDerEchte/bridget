package config

import (
	"errors"
	"os"

	"github.com/nocderechte/bridget/internal/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MQTT struct {
		Topics map[string]byte `yaml:"topics"`
		Host   string          `yaml:"host"`
		Port   int             `yaml:"port"`
	} `yaml:"mqtt"`
	Database struct {
		InfluxDB *struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			Organization string `yaml:"org"`
			Database     string `yaml:"database"`
			AuthScheme   string `yaml:"authScheme"`
		} `yaml:"influx"`
		TailscaleDB *struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			Organization string `yaml:"org"`
			Database     string `yaml:"database"`
			AuthScheme   string `yaml:"authScheme"`
		} `yaml:"tailscale"`
	} `yaml:"database"`
}

func (c *Config) validate() error {
	if c.MQTT.Host == "" {
		return errors.New("missing/invalid config key mqtt.host")
	}
	if c.MQTT.Port <= 0 || c.MQTT.Port > 65535 {
		return errors.New("missing/invalid config key mqtt.port")
	}
	if c.MQTT.Topics == nil {
		return errors.New("missing/invalid config key mqtt.topics")
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	logging.Debug("Opening config file.")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	logging.Debug("Successfully opened config file.")

	defer closeConfig(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	logging.Debug("Decoding contents of config file.")

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	logging.Debug("Successfully decoded contents of config file.")

	logging.Debug("Validating config.")
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

func closeConfig(f *os.File) {
	logging.Debug("Closing config file.")
	if err := f.Close(); err != nil {
		logging.Warn("Failed to close config file.")
	}
	logging.Debug("Successfully closed config file.")
}
