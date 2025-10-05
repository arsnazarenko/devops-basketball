package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

const pathToConfig = "./config/config.yml"

type (
	Config struct {
		HTTP     `yaml:"http"`
		Postgres `yaml:"postgres"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Postgres struct {
		PostgresURL string `env-required:"true" yaml:"postgres_url" env:"POSTGRES_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(pathToConfig, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
