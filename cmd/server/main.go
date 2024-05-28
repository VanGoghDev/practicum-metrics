package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/filestorage"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("failed to run app %w", err)
	}
}

func run() error {
	// config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config %w", err)
	}

	// logger
	zlog, err := logger.New(cfg.Loglevel)
	if err != nil {
		return fmt.Errorf("failed to init logger %w", err)
	}
	zlog.Info("Logger init")
	zlog.Sugar().Info(cfg)

	// storage
	s, err := filestorage.New(zlog, cfg)
	if err != nil {
		return fmt.Errorf("failed to init storage %w", err)
	}

	// router
	router := chirouter.BuildRouter(s, s, zlog)

	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
