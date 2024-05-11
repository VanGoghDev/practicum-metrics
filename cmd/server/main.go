package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
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

	// storage
	s, err := memstorage.New()
	if err != nil {
		return fmt.Errorf("failed to init storage %w", err)
	}

	// router
	router := chirouter.BuildRouter(&s, &s)

	log.Printf("Server start and running on port %s \n", cfg.Address)
	err = http.ListenAndServe(cfg.Address, router)
	return fmt.Errorf("failed to serve: %w", err)
}
