package configs

import (
	"github.com/go-faster/errors"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	EnvPath                  = ".env"
	ErrLoadCfgFile           = "failed to load dot-env file"
	ErrUnmarshalCfgsFromFile = "failed to load configurations"
)

func loadCfg(c *Config) error {
	var errLoad error
	if err := godotenv.Load(EnvPath); err != nil {
		errLoad = errors.Wrap(err, ErrLoadCfgFile)
	}

	if err := envconfig.Process("", c); err != nil {
		if errLoad != nil {
			return errors.Wrap(errors.Wrap(err, ErrUnmarshalCfgsFromFile), errLoad.Error())
		}
		return errors.Wrap(err, ErrUnmarshalCfgsFromFile)
	}

	return nil
}
