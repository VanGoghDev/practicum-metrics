package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/app"
	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	"github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
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

	// storage
	s, err := memstorage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage %w", err)
	}

	sapp, err := app.New(cfg, zlog, &s)
	if err != nil {
		return fmt.Errorf("failed to init app %w", err)
	}

	if cfg.StoreInterval > 0 {
		go func() {
			err := sapp.RunApp()
			zlog.Error(fmt.Sprintf("failed to run app: %v", err))
		}()
	}

	// router
	router := chirouter.BuildRouter(&s, &s, zlog)

	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
