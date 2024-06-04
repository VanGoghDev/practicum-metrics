package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Address            string `env:"ADDRESS"`
	Loglevel           string `env:"LOGLVL"`
	FileStoragePath    string `env:"FILE_STORAGE_PATH"`
	Restore            bool   `env:"RESTORE"`
	StoreInterval      time.Duration
	DBConnectionString string `env:"DATABASE_DSN"`
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
	var flagAddress, flagFileStoragePath, flagLoglevel, flagDBConnection string
	var flagRestore bool
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLoglevel, "lvl", "info", "log level")
	flag.Int64Var(&flagStoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&flagFileStoragePath,
		"f", "/tmp/metrics-db.json", "path to file storage")
	flag.StringVar(&flagDBConnection, "d", "", "db connection string")
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

	if _, present := os.LookupEnv("DATABASE_DSN"); !present {
		cfg.DBConnectionString = flagFileStoragePath
	}

	if _, present := os.LookupEnv("RESTORE"); !present {
		cfg.Restore = flagRestore
	}

	if v, present := os.LookupEnv("STORE_INTERVAL"); !present {
		cfg.StoreInterval = time.Duration(flagStoreInterval) * time.Second
	} else {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("unable to set storeInterval value: %w", err)
		}
		cfg.StoreInterval = time.Duration(i) * time.Second
	}

	return &cfg, nil
}
