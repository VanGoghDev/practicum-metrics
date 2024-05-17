package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORTINTERVAL"`
	PollInterval   time.Duration `env:"POLLINTERVAL"`
}

const (
	defaultReportInterval int64 = 10
	defaultPollInterval   int64 = 2
)

func Load() (config *Config, err error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %w", err)
	}

	var flagAddress string
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")

	var reportInteval, pollInterval int64
	flag.Int64Var(&reportInteval,
		"r", defaultReportInterval,
		"report interval (interval of requests to consumer, in seconds)")
	flag.Int64Var(&pollInterval, "p", defaultPollInterval, "poll interval (interval of metrics fetch, in seconds)")

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
