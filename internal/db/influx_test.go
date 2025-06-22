package db_test

// import (
// 	"os"
// 	"testing"

// 	"github.com/nocderechte/bridget/internal/config"
// 	"github.com/nocderechte/bridget/internal/db"
// 	"github.com/nocderechte/bridget/internal/logging"
// 	"github.com/nocderechte/bridget/internal/mqtt"
// )

// func BenchmarkWriteInfluxDB(b *testing.B) {
// 	b.StopTimer()
// 	logging.InitLogger()

// 	// load config
// 	path := "/home/niklas/Repos/github.com/nocderechte/mqtt-db-connector/internal/config/config.yml"

// 	config, err := config.LoadConfig(path)
// 	logging.Info("loading config.")
// 	if err != nil {
// 		logging.Errorf("Unable to load config.", err)
// 		return
// 	}

// 	// connect to mqtt broker
// 	logging.Info("Establishing connection to mqtt broker.")
// 	if err := mqtt.EstablishConnection(config); err != nil {
// 		logging.Errorf("Failed to establish connection to mqtt broker.", err)
// 		os.Exit(1)
// 	}
// 	defer mqtt.EndConnection()

// 	// init connection to database
// 	logging.Info("Initializing InfluxDB client.")
// 	if err := db.InitInfluxDB(config); err != nil {
// 		logging.Errorf("Failed to initialize InfluxDB client.", err)
// 		os.Exit(1)
// 	}

// 	// test connection to database
// 	logging.Info("Testing database connection.")
// 	if err := db.TestConnection(config); err != nil {
// 		logging.Errorf("Failed to reach database.", err)
// 		os.Exit(1)
// 	}

// 	b.StartTimer()
// 	for range 10000 {
// 		go mqtt.Publish("test-db/benchmark", 0, `{"tags": {"room": "living_room"},"fields": {"temp": "24.2", "hum": 40, "cp": "15i"}}`)
// 	}
// 	b.StopTimer()
// 	b.StartTimer()
// }
