package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"hass/internal/hass"
)

const maxDomains = 5

func newHassClient() (*hass.Client, error) {
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
		return nil, errors.New("no domains defined")
	}

	client, err := hass.New(hass.Config{
		Disabled:         os.Getenv("HASS_DISABLED") == "true",
		Token:            token,
		OverrideEndpoint: endpoint,
		Domains:          domains,
	})
	if err != nil {
		wood.Infof("Credentials path: %v", credentialsPath)
		return nil, err
	}
	return client, nil
}
