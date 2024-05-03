package main

import (
	"fmt"
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	memStorage "github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
)

func main() {
	parseFlags()
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	// config

	// logger

	// storage
	s, err := memStorage.New()
	if err != nil {
		panic("Failed to init storage")
	}

	// router
	router := chirouter.BuildRouter(&s, &s)

	fmt.Printf("Server start and running on port %s \n", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, router)
}
