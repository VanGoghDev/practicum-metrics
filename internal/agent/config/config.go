package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORTINTERVAL"`
	PollInterval   time.Duration `env:"POLLINTERVAL"`
}

func Load() (config *Config, err error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Println("Failed to parse config")
		return nil, err
	}

	var flagAddress string
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")

	var reportInteval, pollInterval int64
	flag.Int64Var(&reportInteval, "r", 10, "report interval (interval of requests to consumer, in seconds)")
	flag.Int64Var(&pollInterval, "p", 2, "poll interval (interval of metrics fetch, in seconds)")

	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = flagAddress
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = time.Duration(reportInteval) * time.Second
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = time.Duration(pollInterval) * time.Second
	}

	return &cfg, nil
}
