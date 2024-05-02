package main

import (
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/routers/chirouter"
	memStorage "github.com/VanGoghDev/practicum-metrics/internal/storage/memstorage"
)

func main() {
	// config

	// logger

	// storage
	s, err := memStorage.New()
	if err != nil {
		panic("Failed to init storage")
	}

	// router
	router := chirouter.BuildRouter(&s, &s)

	http.ListenAndServe(`:8080`, router)
}
