package paho

import (
	"io"
	"log"
	"time"

	"github.com/cwdot/stdlib-go/wood"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

type Connection struct {
	client       MQTT.Client
	disconnected bool
}

func New(config Config) *Connection {
	options := MQTT.NewClientOptions()
	options.AddBroker(config.Broker)
	options.SetClientID(config.ClientID)
	options.SetUsername(config.Username)
	options.SetPassword(config.Password)
	options.SetKeepAlive(time.Minute * 5)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetConnectRetryInterval(time.Second * 5)
	options.SetConnectionLostHandler(func(_ MQTT.Client, err error) {
		wood.Errorf("Connection lost: %v", err)
	})

	lb := loggingBridge{}

	client := MQTT.NewClient(options)
	MQTT.ERROR = log.New(lb, "[ERROR] ", 0)
	MQTT.CRITICAL = log.New(lb, "[CRIT] ", 0)
	//MQTT.WARN = log.New(lb, "[WARN]  ", 0)
	//MQTT.DEBUG = log.New(lb, "[DEBUG] ", 0)
	return &Connection{client: client}
}

func (m *Connection) Connect() error {
	if m.disconnected {
		return errors.New("already disconnected")
	}
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *Connection) Disconnect() {
	m.client.Disconnect(5000)
	m.disconnected = true
}

func (m *Connection) Publish(topic *Topic, payload string, retained bool) error {
	if err := topic.Validate(); err != nil {
		return errors.Wrap(err, "publishing")
	}
	if token := m.client.Publish(topic.String(), 2, retained, payload); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "publishing")
	}

	//wood.Debugf("[%s] Published value", color.Green.It(topic))
	return nil
}

func (m *Connection) Subscribe(topic *Topic, callback Callback) error {
	if err := topic.Validate(); err != nil {
		return errors.Wrap(err, "subscribing")
	}
	cb := func(client MQTT.Client, message MQTT.Message) {
		callback(topic.String(), string(message.Payload()))
	}
	if token := m.client.Subscribe(topic.String(), 1, cb); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

type loggingBridge struct {
	io.Writer
}

func (l loggingBridge) Write(p []byte) (n int, err error) {
	wood.Infof("%s", p)
	return len(p), nil
}

type Callback func(topic string, payload string)
