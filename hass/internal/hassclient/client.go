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

func (c *Client) LightOn(entityId string, opts ...func(*LightOnOpts)) error {
	opt := &LightOnOpts{}
	for _, o := range opts {
		o(opt)
	}

	off := true
	if opt.Brightness > 0 {
		off = false
	}

	arguments := map[string]any{
		"entity_id": entityId,
	}
	if opt.Brightness >= 0 {
		arguments["brightness"] = opt.Brightness
	}

	if opt.Color != nil {
		k, v := opt.Color.Values()
		arguments[k] = v

		// Need to make a call to set the commit the option
		wood.Debugf("Switch options: %s", entityId)
		err := c.Service("light", "turn_on", arguments)
		if err != nil {
			return errors.Wrapf(err, "failed to turn on light: %v", entityId)
		}
		off = false
	}

	if opt.Flash != "" {
		// Needs to the part of the last call to self turn off
		arguments["flash"] = opt.Flash
		off = false
	}

	payload, err := json.Marshal(arguments)
	if err != nil {
		return errors.Wrapf(err, "marshal arguments: %v", arguments)
	}

	action := "on"
	if off {
		action = "off"
	}
	wood.Infof("Turning %s light: %s == %v", action, entityId, string(payload))

	err = c.Service("light", "turn_on", arguments)
	if err != nil {
		return errors.Wrapf(err, "failed to turn on light: %v", entityId)
	}

	if opt.TurnOff != 0 {
		return c.Deactivate(entityId, opt.TurnOff)
	}

	return nil
}

func (c *Client) LightOff(entityId string) error {
	err := c.ServiceSimple("light", "turn_off", entityId)
	if err != nil {
		return errors.Wrapf(err, "failed to turn off light: %v", entityId)
	}
	return nil
}

func (c *Client) Deactivate(entityId string, duration time.Duration) error {
	time.Sleep(duration)
	err := c.LightOff(entityId)
	if err != nil {
		return errors.Wrapf(err, "failed to turn off light for pseudo-transition: %v", entityId)
	}
	return nil
}

func (c *Client) Service(domain string, service string, arguments map[string]any) error {
	err := c.post(fmt.Sprintf("api/services/%s/%s", domain, service), arguments)
	if err != nil {
		return errors.Wrapf(err, "failed to call service %s.%s with %v", domain, service, arguments)
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
			return errors.Wrap(err, "failed to post api")
		}

		text, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read body")
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
