package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/logging"
)

var (
	token  = os.Getenv("INFLUXDB_TOKEN")
	client *influxdb3.Client
)

type InfluxDB struct {
	config *config.Config
}

func (db *InfluxDB) Init() error {
	url := fmt.Sprintf("http://%s:%d/", db.config.Database.InfluxDB.Host, db.config.Database.InfluxDB.Port)
	clientConfig := influxdb3.ClientConfig{
		Host:         url,
		Token:        token,
		Database:     db.config.Database.InfluxDB.Database,
		AuthScheme:   db.config.Database.InfluxDB.AuthScheme,
		Organization: db.config.Database.InfluxDB.Organization,
	}

	var err error
	client, err = influxdb3.New(clientConfig)
	if err != nil {
		return err
	}
	logging.Debug("Successfully validated InfluxDB connection config")
	return nil
}

func (db *InfluxDB) TestConnection() error {
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

func (db *InfluxDB) ProcessWrite(mqttTopic string, payload string) {
	if client == nil {
		logging.Error("InfluxDB client is not initialized.")
	}
	parts := strings.Split(mqttTopic, "/")
	if len(parts) < 2 {
		logging.Error("Invalid topic format - %s", mqttTopic)
	}

	var values influxdb3.PointValues
	if err := json.Unmarshal([]byte(payload), &values); err != nil {
		logging.Errorf("Invalid payload format - %w", err)
	}

	point := influxdb3.NewPoint(parts[1], values.Tags, values.Fields, time.Now())

	logging.Debug("Writing to InfluxDB.")
	err := client.WritePoints(context.Background(), []*influxdb3.Point{point}, influxdb3.WithPrecision(lineprotocol.Nanosecond))
	if err != nil {
		logging.Warn("Failed writing InfluxDB.")
	}
	logging.Debug("Successfully send data to InfluxDB")
}

func (db *InfluxDB) CloseConnection() error {
	if err := client.Close(); err != nil {
		logging.Debug("Failed to gracefully close InfluxDB connection.")
		return err
	}

	logging.Debug("Successfully closed InfluxDB connection.")
	return nil
}
