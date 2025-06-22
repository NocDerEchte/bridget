package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/db"
	"github.com/nocderechte/bridget/internal/logging"
	"github.com/nocderechte/bridget/internal/mqtt"
)

func main() {
	logging.InitLogger()

	// load config
	path := "/home/niklas/Repos/github.com/nocderechte/mqtt-db-connector/internal/config/config.yml"

	config, err := config.LoadConfig(path)
	logging.Info("loading config.")
	if err != nil {
		logging.Errorf("Unable to load config.", err)
		return
	}

	// connect to mqtt broker
	logging.Info("Establishing connection to mqtt broker.")
	if err := mqtt.EstablishConnection(config); err != nil {
		logging.Errorf("Failed to establish connection to mqtt broker.", err)
		os.Exit(1)
	}
	defer mqtt.EndConnection()

	// init connection to database
	logging.Info("Initializing database.")
	if err := db.InitDB(config); err != nil {
		logging.Errorf("Failed to initialize database.", err)
		os.Exit(1)
	}

	// test connection to database
	logging.Info("Testing database connection.")
	if err := db.TestConnection(config); err != nil {
		logging.Errorf("Failed to reach database.", err)
		os.Exit(1)
	}

	logging.Info("Successfully initialized all required components. Ready for requests.")
	// run until Ctrl+C or SIGTERM is send
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
}
