package clientfactory

import (
	"github.com/cwdot/stdlib-go/wood"
	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func LoadAndValidateEnv(credentialsPath string, keys []string) (map[string]string, error) {
	env, err := Load(credentialsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "env loading")
	}

	if err := Validate(env, keys); err != nil {
		wood.Infof("Credentials path: %v", credentialsPath)
		return nil, errors.Wrapf(err, "env validation")
	}

	return env, nil
}

func Load(credentialsPath string) (map[string]string, error) {
	env, err := godotenv.Read(credentialsPath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find %s", credentialsPath)
	}
	return env, nil
}

func Validate(env map[string]string, keys []string) error {
	var errs *multierror.Error
	for _, key := range keys {
		if _, ok := env[key]; !ok {
			errs = multierror.Append(errs, errors.Errorf("missing %s", key))
		}
	}
	return errs.ErrorOrNil()
}
