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

	"github.com/cwdot/stdlib-go/wood"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const maxDomains = 5

var timeout = 3 * time.Second

func New(overrideEndpoint string) (*Client, error) {
	disabled := os.Getenv("HASS_DISABLED")
	if disabled != "" {
		wood.Infof("HASS_DISABLED env var set; exiting early")
		return nil, errors.New("hass disabled")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "error finding home dir")
	}

	// aka ~/.config/hass/credentials.env
	credentialsPath := filepath.Join(home, ".config", "hass", "credentials.env")
	env, err := godotenv.Read(credentialsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find %s", credentialsPath)
	}

	token, ok := env["HASS_TOKEN"]
	if !ok {
		wood.Infof("Credentials path: %v", credentialsPath)
		return nil, errors.New("failed to find hass token")
	}

	domains := make([]string, 0, maxDomains)
	for i := 0; i < maxDomains; i++ {
		value, ok := env[fmt.Sprintf("DOMAIN%d", i)]
		if ok {
			domains = append(domains, value)
		}
	}
	if len(domains) == 0 {
		wood.Infof("Credentials path: %v", credentialsPath)
		return nil, errors.New("no domains defined")
	}

	return &Client{domains: domains, token: token, overrideEndpoint: overrideEndpoint}, nil
}

type Client struct {
	domains          []string
	token            string
	overrideEndpoint string
}

func (c *Client) Execute(entityId string, opts ...func(*LightOnOpts)) error {
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
		req.Header.Set("Authorization", "Bearer "+c.token)

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

	if c.overrideEndpoint != "" {
		return invoke(c.overrideEndpoint)
	}

	var err error
	for _, domain := range c.domains {
		if err = invoke(domain); err == nil {
			break
		}
	}
	return err
}
