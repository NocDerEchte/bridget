package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

func processInfluxDB(mqttTopic string, payload string) error {
	parts := strings.Split(mqttTopic, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid topic format - %s", mqttTopic)
	}

	var values influxdb3.PointValues
	if err := json.Unmarshal([]byte(payload), &values); err != nil {
		return fmt.Errorf("invalid payload format - %w", err)
	}

	point := influxdb3.NewPoint(parts[1], values.Tags, values.Fields, time.Now())
	writeInfluxDB(point)

	return nil
}
