package config

import (
	"errors"
	"os"
	"regexp"

	"github.com/nocderechte/bridget/internal/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	log  logging.Logger
	MQTT struct {
		Topics              map[string]byte `yaml:"topics"`
		Host                string          `yaml:"host"`
		Port                string          `yaml:"port"`
		TimeoutMilliseconds uint            `yaml:"timeoutMilliseconds"`
	} `yaml:"mqtt"`
	Database struct {
		InfluxDB *struct {
			Host           string `yaml:"host"`
			Port           string `yaml:"port"`
			Organization   string `yaml:"org"`
			Database       string `yaml:"database"`
			AuthScheme     string `yaml:"authScheme"`
			TimeoutSeconds int64  `yaml:"timeoutSeconds"`
		} `yaml:"influx"`
		TailscaleDB *struct {
			Host           string `yaml:"host"`
			Port           int    `yaml:"port"`
			Organization   string `yaml:"org"`
			Database       string `yaml:"database"`
			AuthScheme     string `yaml:"authScheme"`
			TimeoutSeconds int    `yaml:"timeoutSeconds"`
		} `yaml:"tailscale"`
	} `yaml:"database"`
}

func New(l logging.Logger) *Config {
	return &Config{log: l}
}

func (c *Config) validate() error {
	portRegex := `^((6553[0-5])|(655[0-2][0-9])|(65[0-4][0-9][0-9])|(6[0-4][0-9][0-9][0-9])|([1-5][0-9][0-9][0-9][0-9])|(\d{1,4}))$`

	if c.MQTT.Host == "" {
		return errors.New("value of config key mqtt.host must not be empty")
	}

	if c.MQTT.Port == "" {
		return errors.New("value of config key mqtt.port must not be empty")
	}

	ok, _ := regexp.MatchString(portRegex, c.MQTT.Port)
	if !ok {
		return errors.New("value of config key mqtt.port is not a valid port number")
	}

	if c.MQTT.Topics == nil {
		return errors.New("value of config key mqtt.topics must not be empty")
	}

	return nil
}

func (c *Config) Load(path string) (*Config, error) {
	c.log.Debug("Opening config file.")
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	c.log.Debug("Successfully opened config file.")

	defer c.close(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	c.log.Debug("Decoding contents of config file.")

	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	c.log.Debug("Successfully decoded contents of config file.")

	c.log.Debug("Validating config.")
	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) close(f *os.File) {
	c.log.Debug("Closing config file.")
	err := f.Close()
	if err != nil {
		c.log.Warn("Failed to close config file.")
	}
	c.log.Debug("Successfully closed config file.")
}
