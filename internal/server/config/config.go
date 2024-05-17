package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func Load() (config *Config, err error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables %w", err)
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()
	return &cfg, nil
}
