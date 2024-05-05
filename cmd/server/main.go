package main

import (
	"fmt"
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
	cfg := config.MustLoad()
	// logger

	// storage
	s, err := memStorage.New()
	if err != nil {
		panic("Failed to init storage")
	}

	// router
	router := chirouter.BuildRouter(&s, &s)

	fmt.Printf("Server start and running on port %s \n", cfg.Address)
	return http.ListenAndServe(cfg.Address, router)
}
