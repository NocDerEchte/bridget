package db

import (
	"context"
	"fmt"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/logging"
	"github.com/nocderechte/bridget/pkg/helper"
)

var (
	token  = helper.GetEnv("INFLUXDB_TOKEN", "apiv3_9xA-23zqAYoF8FXs7tqCXNYauh2aHXhzpkrQxIoO7HmNGo25EBhtmXkiVNP_bKeHI8OsqoEW7V4Wr1g7llvyZQ")
	client *influxdb3.Client
)

func initInfluxDB(c *config.Config) error {

	url := fmt.Sprintf("http://%s:%d/", c.Database.InfluxDB.Host, c.Database.InfluxDB.Port)
	clientConfig := influxdb3.ClientConfig{
		Host:         url,
		Token:        token,
		Database:     c.Database.InfluxDB.Database,
		AuthScheme:   c.Database.InfluxDB.AuthScheme,
		Organization: c.Database.InfluxDB.Organization,
	}

	var err error
	client, err = influxdb3.New(clientConfig)
	if err != nil {
		return err
	}
	logging.Debug("Successfully validated InfluxDB connection config")
	return nil
}

func testInfluxDBConnection() error {
	if client == nil {
		return fmt.Errorf("influxDB client is not initialized")
	}

	query := `SELECT 1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	if _, err := client.Query(ctx, query); err != nil {
		return err
	}

	return nil
}

func writeInfluxDB(p *influxdb3.Point) error {
	if client == nil {
		return fmt.Errorf("influxDB client is not initialized")
	}

	logging.Debug("Writing to InfluxDB.")
	if err := client.WritePoints(context.Background(), []*influxdb3.Point{p}, influxdb3.WithPrecision(lineprotocol.Nanosecond)); err != nil {
		logging.Warn("Failed writing InfluxDB.")
		return nil
	}
	logging.Debug("Successfully send data to InfluxDB")
	return nil
}

func closeInfluxDBConnection() error {
	if err := client.Close(); err != nil {
		logging.Debug("Failed to close InfluxDB connection")
		return err
	}

	logging.Debug("Successfully closed InfluxDB connection")
	return nil
}
