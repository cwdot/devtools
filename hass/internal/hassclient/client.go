package hassclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
)

var timeout = 3 * time.Second

type Config struct {
	Disabled         bool
	Token            string
	OverrideEndpoint string
	Domains          []string
}

func New(config Config) (*Client, error) {
	if config.Disabled {
		wood.Infof("HASS disabled; exiting early")
		return nil, errors.New("hass disabled")
	}
	return &Client{
		config: config,
	}, nil
}

type Client struct {
	config Config
}

func (c *Client) LightOn(entityId string) error {
	err := c.ServiceSimple("light", "turn_on", entityId)
	if err != nil {
		return errors.Wrapf(err, "turn on light: %v", entityId)
	}
	wood.Infof("Simple light on: %s", entityId)
	return nil
}

func (c *Client) LightOff(entityId string) error {
	err := c.ServiceSimple("light", "turn_off", entityId)
	if err != nil {
		return errors.Wrapf(err, "turn off light: %v", entityId)
	}
	wood.Infof("Simple light off: %s", entityId)
	return nil
}

func (c *Client) Execute(domain string, service string, entityId string, arguments map[string]any) error {
	arguments["entity_id"] = entityId

	payload, err := json.Marshal(arguments)
	if err != nil {
		return errors.Wrapf(err, "marshal arguments: %v", arguments)
	}

	wood.Infof("Turning %s %s: %s == %v", domain, entityId, string(payload))

	err = c.Service(domain, service, arguments)
	if err != nil {
		return errors.Wrapf(err, "turn on light: %v", entityId)
	}

	return nil
}

func (c *Client) Deactivate(entityId string, duration time.Duration) error {
	time.Sleep(duration)
	err := c.LightOff(entityId)
	if err != nil {
		return errors.Wrapf(err, "turn off light for pseudo-transition: %v", entityId)
	}
	return nil
}

func (c *Client) Service(domain string, service string, arguments map[string]any) error {
	err := c.post(fmt.Sprintf("api/services/%s/%s", domain, service), arguments)
	if err != nil {
		return errors.Wrapf(err, "call service %s.%s with %v", domain, service, arguments)
	}
	return nil
}

func (c *Client) ServiceSimple(domain string, service string, entityId string) error {
	arguments := map[string]any{
		"entity_id": entityId,
	}

	return c.Service(domain, service, arguments)
}

func (c *Client) post(endpoint string, arguments map[string]any) error {
	postBody, _ := json.Marshal(arguments)
	requestBody := bytes.NewBuffer(postBody)

	payload, _ := json.Marshal(arguments)
	wood.Tracef("Invoked %s with: %s", endpoint, string(payload))

	client := http.Client{
		Timeout: timeout,
	}

	invoke := func(domain string) error {
		url := fmt.Sprintf("%s/%s", domain, endpoint)
		req, err := http.NewRequest(http.MethodPost, url, requestBody)
		if err != nil {
			return err
		}

		wood.Debugf("Invoked POST %s", url)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.config.Token)

		res, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "post api")
		}

		text, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "read body")
		}
		if res.StatusCode != 200 {
			return errors.Errorf("unknown code: %d => %s", res.StatusCode, string(text))
		}
		return nil
	}

	if c.config.OverrideEndpoint != "" {
		return invoke(c.config.OverrideEndpoint)
	}

	var err error
	for _, domain := range c.config.Domains {
		if err = invoke(domain); err == nil {
			break
		}
	}
	return err
}
