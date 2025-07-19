package db

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/logging"
)

var ()

type InfluxDB struct {
	config *config.Config
	log    *logging.Logger
	client *influxdb3.Client
	token  string
}

func (db *InfluxDB) Init() error {
	url := "http://" + net.JoinHostPort(db.config.Database.InfluxDB.Host, db.config.Database.InfluxDB.Port)
	db.token = os.Getenv("INFLUXDB_TOKEN")
	clientConfig := influxdb3.ClientConfig{
		Host:         url,
		Token:        db.token,
		Database:     db.config.Database.InfluxDB.Database,
		AuthScheme:   db.config.Database.InfluxDB.AuthScheme,
		Organization: db.config.Database.InfluxDB.Organization,
	}

	var err error
	db.client, err = influxdb3.New(clientConfig)
	if err != nil {
		return err
	}
	db.log.Debug("Successfully validated InfluxDB connection config")
	return nil
}

func (db *InfluxDB) TestConnection() error {
	if db.client == nil {
		return errors.New("influxDB client is not initialized")
	}

	query := `SELECT 1`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*time.Duration(db.config.Database.InfluxDB.TimeoutSeconds),
	)
	defer cancel()

	if _, err := db.client.Query(ctx, query); err != nil {
		return err
	}

	return nil
}

func (db *InfluxDB) ProcessWrite(mqttTopic string, payload string) {
	if db.client == nil {
		db.log.Error("InfluxDB client is not initialized.")
	}
	parts := strings.Split(mqttTopic, "/")
	requiredParts := 2
	if len(parts) < requiredParts {
		db.log.Error("Invalid topic format - %s", mqttTopic)
	}

	var values influxdb3.PointValues
	if err := json.Unmarshal([]byte(payload), &values); err != nil {
		db.log.Errorf("Invalid payload format - %w", err)
	}

	point := influxdb3.NewPoint(parts[1], values.Tags, values.Fields, time.Now())

	db.log.Debug("Writing to InfluxDB.")
	err := db.client.WritePoints(
		context.Background(),
		[]*influxdb3.Point{point},
		influxdb3.WithPrecision(lineprotocol.Nanosecond),
	)
	if err != nil {
		db.log.Warn("Failed writing InfluxDB.")
	}
	db.log.Debug("Successfully send data to InfluxDB")
}

func (db *InfluxDB) CloseConnection() error {
	if err := db.client.Close(); err != nil {
		db.log.Debug("Failed to gracefully close InfluxDB connection.")
		return err
	}

	db.log.Debug("Successfully closed InfluxDB connection.")
	return nil
}
