package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UserTokenExpiration   time.Duration `envconfig:"USER_TOKEN_EXPIRATION" default:"720h"`
	ServiceAllowedOrigins []string      `envconfig:"SERVICE_ALLOWED_ORIGINS"`
}

func New() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
