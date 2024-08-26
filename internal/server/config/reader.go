package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
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
		return fmt.Errorf("%w: failed to parse environment variables", err)
	}
	return nil
}

type FlagReader struct {
}

func (r *FlagReader) ReadFlags(cfg *Config) error {
	if cfg == nil {
		return errors.New("config is nil")
	}

	var flagStoreInterval int64
	var flagAddress, flagFileStoragePath, flagLoglevel, flagDBConnection, flagKey string
	var flagRestore bool
	flag.StringVar(&flagAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLoglevel, "lvl", "info", "log level")
	flag.StringVar(&flagKey, "k", "", "signature key")
	flag.Int64Var(&flagStoreInterval, "i", defaultStoreInterval, "store interval in seconds")
	flag.StringVar(&flagFileStoragePath, "f", "", "path to file storage")
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
		cfg.DBConnectionString = flagDBConnection
	}

	if _, present := os.LookupEnv("RESTORE"); !present {
		cfg.Restore = flagRestore
	}

	if v, present := os.LookupEnv("STORE_INTERVAL"); !present {
		cfg.StoreInterval = time.Duration(flagStoreInterval) * time.Second
	} else {
		i, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("unable to set storeInterval value: %w", err)
		}
		cfg.StoreInterval = time.Duration(i) * time.Second
	}

	if _, present := os.LookupEnv("KEY"); !present {
		cfg.Key = flagKey
	}
	return nil
}
