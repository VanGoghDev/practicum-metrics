package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

// Config хранит параметры приложения.
type Config struct {
	// Address адрес на котором запускается приложение.
	Address string `env:"ADDRESS"`

	// Loglevel уровень логирования (DEBUG/ERROR/...).
	Loglevel string `env:"LOGLVL"`

	// FileStoragePath путь до хранения данных в файле.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	// DbConnectionString строка подключения к бд.
	DBConnectionString string `env:"DATABASE_DSN"`

	// Key ключ для шифрования запросов.
	Key string `env:"KEY"`

	// Restore флаг того что нужно восстанавливать данные из файла или нет.
	Restore bool `env:"RESTORE"`

	// StoreInterval интервал с которым сохраняются данные в файловое хранилище.
	StoreInterval time.Duration
}

const (
	defaultStoreInterval int64 = 300
)

// Load инициализирует конфиг.
func Load() (config *Config, err error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables %w", err)
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
			return nil, fmt.Errorf("unable to set storeInterval value: %w", err)
		}
		cfg.StoreInterval = time.Duration(i) * time.Second
	}

	if _, present := os.LookupEnv("KEY"); !present {
		cfg.Key = flagKey
	}

	return &cfg, nil
}
