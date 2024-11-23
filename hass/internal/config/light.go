package config

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/pkg/errors"

	"hass/internal/hassclient"
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

func createLightOpts(entity Light) ([]func(opts *hassclient.LightOnOpts), error) {
	opts := make([]func(opts *hassclient.LightOnOpts), 0, 5)

	switch entity.Color {
	case "red":
		opts = append(opts, hassclient.Red())
	case "green":
		opts = append(opts, hassclient.Green())
	case "blue":
		opts = append(opts, hassclient.Blue())
	case "white":
		opts = append(opts, hassclient.White())
	case "yellow":
		opts = append(opts, hassclient.Yellow())
	case "":
	default:
		return nil, errors.Errorf("unknown color: %s", entity.Color)
	}

	if entity.Brightness >= 0 {
		opts = append(opts, hassclient.Brightness(entity.Brightness))
	}

	switch entity.Flash {
	case "long":
		opts = append(opts, hassclient.LongFlash())
	case "short":
		opts = append(opts, hassclient.ShortFlash())
	case "":
	default:
		wood.Debugf("unknown flash: %s", entity.Flash)
	}

	return opts, nil
}
