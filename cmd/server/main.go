package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/config"
	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	memStorage "github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// logger

	// storage
	s, err := memStorage.New()
	if err != nil {
		log.Fatal(err)
	}

	// router
	router := chirouter.BuildRouter(&s, &s)

	log.Printf("Server start and running on port %s \n", cfg.Address)
	err = http.ListenAndServe(cfg.Address, router)
	return fmt.Errorf("failed to serve: %w", err)
}
