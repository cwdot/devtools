package lightmanager

import (
	"encoding/json"
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"
	"hass/internal/config"
	"hass/internal/hassclient"
)

func New(aliases map[string]string) *LightManager {
	return &LightManager{aliases: aliases}
}

type LightManager struct {
	lights  map[string]config.LightEntity
	aliases map[string]string
}

func (c *LightManager) List() []string {
	keys := make([]string, 0, len(c.aliases))
	for k, _ := range c.aliases {
		keys = append(keys, k)
	}
	return keys
}

func (c *LightManager) GetLightId(alias string) string {
	if fullId, ok := c.aliases[alias]; ok {
		return fullId
	}
	return alias
}

func (c *LightManager) LightOn(client *hassclient.Client, entity config.LightEntity) error {
	id := entity.Light
	if fullId := c.GetLightId(entity.Light); fullId != "" {
		id = fullId
	}

	if err := client.LightOn(id); err != nil {
		return err
	}
	return nil
}

func (c *LightManager) LightOff(client *hassclient.Client, entity config.LightEntity) error {
	id := entity.Light
	if fullId := c.GetLightId(entity.Light); fullId != "" {
		id = fullId
	}

	if err := client.LightOff(id); err != nil {
		return err
	}
	return nil
}

func (c *LightManager) Execute(hc *hassclient.Client, entity config.LightEntity) error {
	opts, err := createLightOpts(entity)
	if err != nil {
		return err
	}

	entityId := entity.Light
	if fullId := c.GetLightId(entity.Light); fullId != "" {
		entityId = fullId
	}

	off := true
	if opts.Brightness > 0 {
		off = false
	}

	arguments := map[string]any{
		"entity_id": entityId,
	}
	if opts.Brightness >= 0 {
		arguments["brightness"] = opts.Brightness
	}

	if opts.Color != nil {
		k, v := opts.Color.Values()
		arguments[k] = v

		// Need to make a call to set the commit the option
		wood.Debugf("Switch options: %s", entityId)
		err := hc.Service("light", "turn_on", arguments)
		if err != nil {
			return errors.Wrapf(err, "turn on light: %v", entityId)
		}
		off = false
	}

	if opts.Flash != "" {
		// Needs to the part of the last call to self turn off
		arguments["flash"] = opts.Flash
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

	err = hc.Service("light", "turn_on", arguments)
	if err != nil {
		return errors.Wrapf(err, "turn on light: %v", entityId)
	}

	if opts.TurnOff != 0 {
		return hc.Deactivate(entityId, opts.TurnOff)
	}

	return nil
}

func createLightOpts(entity config.LightEntity) (*hassclient.LightOnOpts, error) {
	options := make([]func(opts *hassclient.LightOnOpts), 0, 5)

	switch entity.Color {
	case "red":
		options = append(options, hassclient.Red())
	case "green":
		options = append(options, hassclient.Green())
	case "blue":
		options = append(options, hassclient.Blue())
	case "white":
		options = append(options, hassclient.White())
	case "yellow":
		options = append(options, hassclient.Yellow())
	case "":
	default:
		return nil, errors.Errorf("unknown color: %s", entity.Color)
	}

	if entity.Brightness >= 0 {
		options = append(options, hassclient.Brightness(entity.Brightness))
	}

	switch entity.Flash {
	case "long":
		options = append(options, hassclient.LongFlash())
	case "short":
		options = append(options, hassclient.ShortFlash())
	case "":
	default:
		wood.Debugf("unknown flash: %s", entity.Flash)
	}

	opts := &hassclient.LightOnOpts{}
	for _, o := range options {
		o(opts)
	}
	return opts, nil
}
