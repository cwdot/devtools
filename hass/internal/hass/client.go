package hass

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/cwdot/go-stdlib/wood"
)

const (
	domain = "http://192.168.1.101:8123"
	//domain           = "https://quakequack.duckdns.org"
)

func New() (*Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}

	configLocation := filepath.Join(home, ".credentials.env")
	env, err := godotenv.Read(configLocation)
	if err != nil {
		return nil, err
	}

	token, ok := env["HASS_TOKEN"]
	if !ok {
		return nil, errors.New("failed to find hass token")
	}

	return &Client{token: token}, nil
}

type Client struct {
	token string
}

func (c *Client) LightOn(entityId string, opts ...func(*LightOnOpts)) error {
	opt := &LightOnOpts{}
	for _, o := range opts {
		o(opt)
	}

	arguments := map[string]any{
		"entity_id": entityId,
	}

	if opt.Color != nil {
		k, v := opt.Color.Values()
		arguments[k] = v

		// Need to make a call to set the commit the option
		wood.Infof("Switch options: %s", entityId)
		err := c.Service("light", "turn_on", arguments)
		if err != nil {
			return errors.Wrapf(err, "failed to turn on light: %v", entityId)
		}
	}

	if opt.Brightness != 0 {
		arguments["brightness"] = opt.Brightness
	}

	if opt.Flash != "" {
		// Needs to the part of the last call to self turn off
		arguments["flash"] = opt.Flash
	}

	wood.Infof("Turning on light: %s == %v", entityId, arguments)

	err := c.Service("light", "turn_on", arguments)
	if err != nil {
		return errors.Wrapf(err, "failed to turn on light: %v", entityId)
	}

	if opt.TurnOff != 0 {
		time.Sleep(opt.TurnOff)
		err = c.LightOff(entityId)
		if err != nil {
			return errors.Wrapf(err, "failed to turn off light for pseudo-transition: %v", entityId)
		}
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

	wood.Debugf("Invoked %s with: %s", endpoint, string(postBody))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", domain, endpoint), requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to post api")
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}

	return nil
}
