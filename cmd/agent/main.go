package main

import (
	"log"

	agent "github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
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

	zlog.Info("Logger init")

	app := agent.New(zlog, cfg)
	err = app.RunApp()
	if err != nil {
		log.Fatalf("failed to run an app: %v", err)
	}
}
