package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func MustLoad() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Println("Failed to parse environment variables")
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	}
	flag.Parse()
	return &cfg
}
