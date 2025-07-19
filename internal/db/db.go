package db

import (
	"errors"

	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/logging"
)

type Database interface {
	TestConnection() error
	Init() error
	ProcessWrite(mqttTopic string, payload string)
	CloseConnection() error
}

func NewDatabaseFromConfig(c *config.Config, l logging.Logger) (Database, error) {
	switch {
	case c.Database.InfluxDB != nil:
		return &InfluxDB{config: c, log: &l}, nil
	case c.Database.TailscaleDB != nil:
		return nil, errors.New("tailscale is not supported")
	default:
		return nil, errors.New("database configuration is empty")
	}
}
