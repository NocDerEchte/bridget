package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/db"
	"github.com/nocderechte/bridget/internal/logging"
	"github.com/nocderechte/bridget/internal/mqtt"
	"github.com/nocderechte/bridget/pkg/helper"
)

func main() {
	logging.InitLogger()

	// load config
	configPath := helper.GetEnv("BRIDGET_CONFIG", "/etc/bridget/config.yml")

	config, err := config.LoadConfig(configPath)
	logging.Info("loading config.")
	if err != nil {
		logging.Errorf("Unable to load config.", err)
		return
	}

	// init connection to database
	logging.Info("Initializing database.")
	database, err := db.NewDatabaseFromConfig(config)
	if err != nil {
		logging.Errorf("Failed to create database from config.", err)
	}

	err = database.Init()
	if err != nil {
		logging.Errorf("Failed to initialize database.", err)
		os.Exit(1)
	}

	// test connection to database
	logging.Info("Testing database connection.")
	if err := database.TestConnection(); err != nil {
		logging.Errorf("Failed to reach database.", err)
		os.Exit(1)
	}

	// connect to mqtt broker
	logging.Info("Establishing connection to mqtt broker.")
	if err := mqtt.EstablishConnection(config, database); err != nil {
		logging.Errorf("Failed to establish connection to mqtt broker.", err)
		os.Exit(1)
	}
	defer mqtt.EndConnection()

	logging.Info("Successfully initialized all required components. Ready for requests.")
	// run until Ctrl+C or SIGTERM is send
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
}
