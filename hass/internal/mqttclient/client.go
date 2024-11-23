package mqttclient

import (
	"hass/internal/mqttclient/internal/paho"
)

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

func New(config Config) (*Client, error) {
	return &Client{
		config: config,
		mqtt: paho.New(paho.Config{
			Broker:   config.Broker,
			ClientID: config.ClientID,
			Username: config.Username,
			Password: config.Password,
		}),
	}, nil
}

type Client struct {
	config Config
	mqtt   *paho.Connection
}

func (c *Client) Connect() error {
	return c.mqtt.Connect()
}

func (c *Client) Disconnect() {
	c.mqtt.Disconnect()
}

func (c *Client) Publish(topic string, payload string, retained bool) error {
	t := paho.NewTopic(topic, map[string]string{})
	return c.mqtt.Publish(t, payload, retained)
}

type Callback func(topic string, payload string)

func (c *Client) Subscribe(topic string, callback Callback) error {
	t := paho.NewTopic(topic, map[string]string{})
	return c.mqtt.Subscribe(t, func(topic string, payload string) {
		callback(topic, payload)
	})
}
