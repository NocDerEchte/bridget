package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nocderechte/bridget/internal/config"
	"github.com/nocderechte/bridget/internal/db"
	"github.com/nocderechte/bridget/internal/logging"
	"github.com/nocderechte/bridget/pkg/helper"
)

var (
	client             mqtt.Client
	messagePubHandler  mqtt.MessageHandler
	connectHandler     mqtt.OnConnectHandler
	connectLostHandler mqtt.ConnectionLostHandler
)

func EstablishConnection(c *config.Config, db db.Database) error {
	messagePubHandler = func(client mqtt.Client, msg mqtt.Message) {
		go db.ProcessWrite(msg.Topic(), string(msg.Payload()))
	}

	connectHandler = func(client mqtt.Client) {
		logging.Info("Successfully connected to mqtt broker")
		client.SubscribeMultiple(c.MQTT.Topics, messagePubHandler)
	}

	connectLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}

	broker := c.MQTT.Host
	port := c.MQTT.Port
	username := helper.GetEnv("MQTT_USER", "exampleuser")
	password := helper.GetEnv("MQTT_PASSWORD", "examplepassword")
	clientID := helper.GetEnv("MQTT_CLIENTID", "connector")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetPassword(password)

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("unable to connect to mqtt broker: %w", token.Error())
	}

	return nil
}

func EndConnection() {
	client.Disconnect(10)
}

func Publish(topic string, qos byte, payload any) {
	client.Publish(topic, qos, false, payload)
}
