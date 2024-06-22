package main

import (
	"log"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// logger
	zlog, err := logger.New(cfg.Loglevel)
	if err != nil {
		log.Fatal("failed to init logger %w", err)
	}
	zlog.Sugar().Infow("config",
		"rateLimit", cfg.RateLimit,
		"pollInterval", cfg.PollInterval,
		"reportInterval", cfg.ReportInterval,
	)

	zlog.Info("Logger init")

	agent := app.New(zlog, cfg)
	err = agent.RunApp()
	if err != nil {
		log.Fatalf("failed to run an app: %v", err)
	}
}
