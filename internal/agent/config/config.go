package config

import (
	"fmt"
	"time"
)

// Config хранит параметры приложения.
type Config struct {
	// Address адрес сервиса куда отправлять метрики.
	Address string `env:"ADDRESS"`

	// Loglevel уровень логирования (DEBUG/ERROR/...).
	Loglevel string `env:"LOGLVL"`

	// Key ключ подписи.
	Key string `env:"KEY"`

	// RateLimit кол-во горутин, которые будут отправлять запросы на сервер.
	RateLimit int64 `env:"RATE_LIMIT"`

	// ReportInterval интервал с которым происходит отправка метрик на сервер.
	ReportInterval time.Duration `env:"REPORTINTERVAL"`

	// PollInterval интервал с которым происходит опрос метрик.
	PollInterval time.Duration `env:"POLLINTERVAL"`
}

const (
	defaultReportInterval int64 = 10
	defaultPollInterval   int64 = 2
)

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
