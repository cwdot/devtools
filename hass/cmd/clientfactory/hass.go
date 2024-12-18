package clientfactory

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cwdot/stdlib-go/wood"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"hass/internal/hassclient"
)

const maxDomains = 5

func NewHassClient(endpoint string) (*hassclient.Client, error) {
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

	if err := Validate(env, []string{"HASS_TOKEN"}); err != nil {
		return nil, errors.Wrapf(err, "env validation")
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

	client, err := hassclient.New(hassclient.Config{
		Disabled:         os.Getenv("HASS_DISABLED") == "true",
		Token:            env["HASS_TOKEN"],
		OverrideEndpoint: endpoint,
		Domains:          domains,
	})
	if err != nil {
		wood.Infof("Credentials path: %v", credentialsPath)
		return nil, err
	}
	return client, nil
}
