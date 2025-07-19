package mqtt

import (
	"fmt"
	"net"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/db"
	"github.com/nocderechte/bridget/internal/logging"
	"github.com/nocderechte/bridget/pkg/helper"
)

type Client struct {
	log                logging.Logger
	config             *config.Config
	database           db.Database
	client             mqtt.Client
	maxReconnects      int
	reconnects         int
	messagePubHandler  mqtt.MessageHandler
	connectHandler     mqtt.OnConnectHandler
	connectLostHandler mqtt.ConnectionLostHandler
}

func New(config *config.Config, db db.Database, l logging.Logger) *Client {
	return &Client{config: config, database: db, log: l}
}

func (c *Client) EstablishConnection() error {
	c.messagePubHandler = func(_ mqtt.Client, msg mqtt.Message) {
		go c.database.ProcessWrite(msg.Topic(), string(msg.Payload()))
	}

	c.connectHandler = func(client mqtt.Client) {
		c.log.Info("Successfully connected to mqtt broker")
		client.SubscribeMultiple(c.config.MQTT.Topics, c.messagePubHandler)
	}

	c.connectLostHandler = func(client mqtt.Client, err error) {
		if c.reconnects < c.maxReconnects {
			c.log.Warn("lost connection to mqtt broker. Retrying in 5 seconds... - %w", err)

			c.reconnects++

			client.Connect()
			return
		}

		c.log.Errorf("lost connection to mqtt broker - %w", err)
	}

	broker := c.config.MQTT.Host
	port := c.config.MQTT.Port
	username := helper.GetEnv("MQTT_USER", "exampleuser")
	password := helper.GetEnv("MQTT_PASSWORD", "examplepassword")
	clientID := helper.GetEnv("MQTT_CLIENTID", "connector")

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + net.JoinHostPort(broker, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(c.messagePubHandler)
	opts.OnConnect = c.connectHandler
	opts.OnConnectionLost = c.connectLostHandler

	c.client = mqtt.NewClient(opts)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("unable to connect to mqtt broker: %w", token.Error())
	}

	return nil
}

func (c *Client) EndConnection() {
	c.client.Disconnect(c.config.MQTT.TimeoutMilliseconds)
}

func (c *Client) Publish(topic string, qos byte, payload any) {
	c.client.Publish(topic, qos, false, payload)
}
