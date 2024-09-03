package config

import (
	"fmt"
	"time"
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

type EnvironmentReader interface {
	ReadEnvironment(cfg *Config) error
}

type FlagsReader interface {
	ReadFlags(cfg *Config) error
}

// Load инициализирует конфиг.
func Load(envReader EnvironmentReader, flagsReader FlagsReader) (config *Config, err error) {
	cfg := &Config{}

	if err := envReader.ReadEnvironment(cfg); err != nil {
		return nil, fmt.Errorf("%w: failed to parse environment variables", err)
	}

	if err := flagsReader.ReadFlags(cfg); err != nil {
		return nil, fmt.Errorf("%w: failed to read flags", err)
	}

	return cfg, nil
}
