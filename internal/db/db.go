package db

import (
	"fmt"

	"github.com/nocderechte/bridget/internal/config"
)

type Database interface {
	TestConnection() error
	Init() error
	ProcessWrite(mqttTopic string, payload string)
	CloseConnection() error
}

func NewDatabaseFromConfig(c *config.Config) (Database, error) {
	switch {
	case c.Database.InfluxDB != nil:
		return &InfluxDB{c}, nil
	case c.Database.TailscaleDB != nil:
		return nil, fmt.Errorf("tailscale is not supported")
	default:
		return nil, fmt.Errorf("database configuration is empty")
	}
}
