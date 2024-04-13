package config

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"

	"hass/internal/hass"
)

type LightManager struct {
	lights  map[string]Light
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

func createLightOpts(entity Light) ([]func(opts *hass.LightOnOpts), error) {
	opts := make([]func(opts *hass.LightOnOpts), 0, 5)

	switch entity.Color {
	case "red":
		opts = append(opts, hass.Red())
	case "green":
		opts = append(opts, hass.Green())
	case "blue":
		opts = append(opts, hass.Blue())
	case "white":
		opts = append(opts, hass.White())
	case "yellow":
		opts = append(opts, hass.Yellow())
	case "":
	default:
		return nil, errors.Errorf("unknown color: %s", entity.Color)
	}

	if entity.Brightness >= 0 {
		opts = append(opts, hass.Brightness(entity.Brightness))
	}

	switch entity.Flash {
	case "long":
		opts = append(opts, hass.LongFlash())
	case "short":
		opts = append(opts, hass.ShortFlash())
	case "":
	default:
		wood.Debugf("unknown flash: %s", entity.Flash)
	}

	return opts, nil
}
