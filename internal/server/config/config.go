package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Address         string        `env:"ADDRESS"`
	Loglevel        string        `env:"LOGLVL"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH"`
	Restore         bool          `env:"RESTORE"`
	StoreInterval   time.Duration `env:"STORE_INTERVAL"`
}

const (
	defaultStoreInterval int64 = 300
)

func Load() (config *Config, err error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables %w", err)
	}

	var flagStoreInterval int64
	var flagAddress, flagFileStoragePath, flagLoglevel string
	var flagRestore bool
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLoglevel, "lvl", "info", "log level")
	flag.Int64Var(&flagStoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "path to file storage")
	flag.BoolVar(&flagRestore, "r", true, "restore previous state or not")
	flag.Parse()

	if _, present := os.LookupEnv("ADDRESS"); !present {
		cfg.Address = flagAddress
	}

	if _, present := os.LookupEnv("LOGLVL"); !present {
		cfg.Loglevel = flagLoglevel
	}

	if _, present := os.LookupEnv("FILE_STORAGE_PATH"); !present {
		cfg.FileStoragePath = flagFileStoragePath
	}

	_, present := os.LookupEnv("RESTORE")
	if !present {
		cfg.Restore = flagRestore
	}

	if _, present := os.LookupEnv("STORE_INTERVAL"); !present {
		cfg.StoreInterval = time.Duration(flagStoreInterval) * time.Second
	}

	return &cfg, nil
}
