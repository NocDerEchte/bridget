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
	log := logging.NewLogger()
	config := config.New(*log)

	// load config
	configPath := helper.GetEnv("BRIDGET_CONFIG", "/etc/bridget/config.yml")

	config, err := config.Load(configPath)
	log.Info("loading config.")
	if err != nil {
		log.Errorf("Unable to load config.", err)
		return
	}

	// init connection to database
	log.Info("Initializing database.")
	database, err := db.NewDatabaseFromConfig(config, *log)
	if err != nil {
		log.Errorf("Failed to create database from config.", err)
	}

	err = database.Init()
	if err != nil {
		log.Errorf("Failed to initialize database.", err)
		os.Exit(1)
	}

	// test connection to database
	log.Info("Testing database connection.")
	err = database.TestConnection()
	if err != nil {
		log.Errorf("Failed to reach database.", err)
		os.Exit(1)
	}

	// connect to mqtt broker
	log.Info("Establishing connection to mqtt broker.")
	mqttClient := mqtt.New(config, database, *log)
	err = mqttClient.EstablishConnection()
	if err != nil {
		log.Errorf("Failed to establish connection to mqtt broker.", err)
		os.Exit(1)
	}
	defer mqttClient.EndConnection()

	log.Info("Successfully initialized all required components. Ready for requests.")
	// run until Ctrl+C or SIGTERM is send
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
}
