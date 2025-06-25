package config

import (
	"errors"
	"os"

	"github.com/nocderechte/bridget/internal/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MQTT     mqttConfig     `yaml:"mqtt"`
	Database DatabaseConfig `yaml:"database"`
}

type mqttConfig struct {
	Topics map[string]byte `yaml:"topics"`
	Host   string          `yaml:"host"`
	Port   int             `yaml:"port"`
}

type DatabaseConfig struct {
	Type        string            `yaml:"type"`
	InfluxDB    InfluxDBConfig    `yaml:"influx"`
	TimescaleDB TimescaleDBConfig `yaml:"timescale"`
}

type InfluxDBConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Organization string `yaml:"org"`
	Database     string `yaml:"database"`
	AuthScheme   string `yaml:"authScheme"`
}

type TimescaleDBConfig struct{}

func (c *Config) validate() error {
	validDatabases := []string{"influx", "timescale"}

	// validate mqtt config
	if c.MQTT.Host == "" {
		return errors.New("missing/invalid config key mqtt.host")
	}
	if c.MQTT.Port <= 0 || c.MQTT.Port > 65535 {
		return errors.New("missing/invalid config key mqtt.port")
	}
	if c.MQTT.Topics == nil {
		return errors.New("missing/invalid config key mqtt.topics")
	}

	// validate database config
	valid := false
	for _, db := range validDatabases {
		if c.Database.Type == db {
			valid = true
		}
	}
	if !valid {
		return errors.New("missing/invalid config key database.type")
	}

	// influxdb
	// timescale

	return nil
}

func LoadConfig(path string) (*Config, error) {
	logging.Debug("Opening config file.")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	logging.Debug("Successfully opened config file.")

	defer func(f *os.File) error {
		logging.Debug("Closing config file.")
		if err := f.Close(); err != nil {
			return err
		}
		logging.Debug("Successfully closed config file")
		return nil
	}(file)

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

	logging.Debug("Checking for database config.")
	if (config.Database == DatabaseConfig{}) {
		return nil, errors.New("invalid config. Missing database configuration")
	}
	logging.Debug("Successfully checked database config. Config valid.")

	return &config, nil
}
