package main

import (
	"log"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app := app.New(cfg)
	err = app.RunApp()
	if err != nil {
		log.Fatal(err)
	}
}
