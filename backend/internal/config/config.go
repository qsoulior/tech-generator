package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UserTokenExpiration   time.Duration `envconfig:"USER_TOKEN_EXPIRATION" default:"720h"`
	ServiceAllowedOrigins []string      `envconfig:"SERVICE_ALLOWED_ORIGINS" required:"true"`
	Ed25519PrivateKeyPath string        `envconfig:"ED25519_PRIVATE_KEY_PATH" required:"true"`
	Ed25519PublicKeyPath  string        `envconfig:"ED25519_PUBLIC_KEY_PATH" required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
