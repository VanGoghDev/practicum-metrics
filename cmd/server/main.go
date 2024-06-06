package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/logger"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	"github.com/VanGoghDev/practicum-metrics/internal/storage"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal("failed to run app %w", err)
	}
}

func run(ctx context.Context) error {
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
	zlog.Sugar().Debug(cfg)

	// storage
	s, err := storage.New(ctx, cfg, zlog)
	defer func() {
		err = s.Close(ctx)
	}()

	if err != nil {
		return fmt.Errorf("failed to init storage %w", err)
	}

	// router
	router := chirouter.BuildRouter(ctx, s, zlog, cfg)

	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return fmt.Errorf("failed to run app: %w", err)
}
