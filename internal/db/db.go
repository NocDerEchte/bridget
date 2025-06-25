package db

import (
	"fmt"

	"github.com/nocderechte/bridget/internal/config"
)

var testFuncs = map[string]func() error{
	"influx": testInfluxDBConnection,
}

var initFuncs = map[string]func(c *config.Config) error{
	"influx": initInfluxDB,
}

var processFuncs = map[string]func(mqttTopic string, payload string) error{
	"influx": processInfluxDB,
}

var closeFuncs = map[string]func() error{
	"influx": closeInfluxDBConnection,
}

func TestConnection(c *config.Config) error {
	testFunc, ok := testFuncs[c.Database.Type]
	if !ok {
		return fmt.Errorf("invalid database type - %s", c.Database.Type)
	}
	if err := testFunc(); err != nil {
		return fmt.Errorf("%s connection failed - %w", c.Database.Type, err)
	}

	return nil
}

func InitDB(c *config.Config) error {
	initFunc, ok := initFuncs[c.Database.Type]
	if !ok {
		return fmt.Errorf("invalid database type - %s", c.Database.Type)
	}
	if err := initFunc(c); err != nil {
		return fmt.Errorf("%s initialization failed - %w", c.Database.Type, err)
	}

	return nil
}

func ProcessWrite(c *config.Config, mqttTopic string, payload string) error {
	processFunc, ok := processFuncs[c.Database.Type]
	if !ok {
		return fmt.Errorf("invalid database type - %s", c.Database.Type)
	}
	if err := processFunc(mqttTopic, payload); err != nil {
		return fmt.Errorf("%s processing failed - %w", c.Database.Type, err)
	}

	return nil
}

func CloseConnection(c *config.Config) error {
	closeFunc, ok := closeFuncs[c.Database.Type]
	if !ok {
		return fmt.Errorf("invalid database type - %s", c.Database.Type)
	}
	if err := closeFunc(); err != nil {
		return fmt.Errorf("%s connection closing failed - %w", c.Database.Type, err)
	}

	return nil
}
