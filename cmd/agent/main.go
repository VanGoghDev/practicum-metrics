package main

import (
	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
)

func main() {

	cfg := config.MustLoad()

	app := app.New(cfg)
	app.MustRun()
}
