package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORTINTERVAL"`
	PollInterval   time.Duration `env:"POLLINTERVAL"`
}

func MustLoad() *Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		panic("failed to parse environment variables")
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	}

	var reportInteval, pollInterval int64
	if cfg.ReportInterval == 0 {
		flag.Int64Var(&reportInteval, "r", 10, "report interval (interval of requests to consumer, in seconds)")
		cfg.ReportInterval = time.Duration(reportInteval) * time.Second
	}

	if cfg.PollInterval == 0 {
		flag.Int64Var(&pollInterval, "p", 2, "poll interval (interval of metrics fetch, in seconds)")
		cfg.PollInterval = time.Duration(pollInterval) * time.Second
	}

	flag.Parse()
	return &cfg
}
