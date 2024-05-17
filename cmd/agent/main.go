package main

import (
	"log"

	agent "github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app := agent.New(cfg)
	err = app.RunApp()
	if err != nil {
		log.Fatalf("failed to run an app: %v", err)
	}
}
