package main

import (
	"net/http"

	"github.com/VanGoghDev/practicum-metrics/internal/server/handlers/update"
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
	mux := setupMux()

	mux.HandleFunc(`/update/`, update.UpdateHandler(&s))

	http.ListenAndServe(`:8080`, mux)
}

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}
