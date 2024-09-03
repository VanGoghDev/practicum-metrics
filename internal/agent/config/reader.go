package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

type EnvironmentReader interface {
	ReadEnvironment(cfg *Config) error
}

type FlagsReader interface {
	ReadFlags(cfg *Config) error
}

type EnvReader struct {
}

func (r *EnvReader) ReadEnvironment(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("failed to parse config %w", err)
	}
	return nil
}

type FlagReader struct {
}

func (r *FlagReader) ReadFlags(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	var reportInteval, pollInterval, rateLimit int64
	var logLevel, flagAddress, flagKey string

	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&reportInteval,
		"r", defaultReportInterval,
		"report interval (interval of requests to consumer, in seconds)")
	flag.Int64Var(&pollInterval, "p", defaultPollInterval, "poll interval (interval of metrics fetch, in seconds)")
	flag.StringVar(&logLevel, "lvl", "info", "log level")
	flag.StringVar(&flagKey, "k", "", "signature key")
	flag.Int64Var(&rateLimit, "l", 1, "number of goroutines for sending metrics to server")

	flag.Parse()

	if _, present := os.LookupEnv("ADDRESS"); !present {
		cfg.Address = flagAddress
	}

	if _, present := os.LookupEnv("LOGLVL"); !present {
		cfg.Loglevel = logLevel
	}

	if _, present := os.LookupEnv("REPORTINTERVAL"); !present {
		cfg.ReportInterval = time.Duration(reportInteval) * time.Second
	}

	if _, present := os.LookupEnv("POLLINTERVAL"); !present {
		cfg.PollInterval = time.Duration(pollInterval) * time.Second
	}

	if _, present := os.LookupEnv("KEY"); !present {
		cfg.Key = flagKey
	}

	if _, present := os.LookupEnv("RATE_LIMIT"); !present {
		cfg.RateLimit = rateLimit
	}

	return nil
}
