package main

import (
	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
)

func main() {
	app := app.New()
	app.MustRun()
}
