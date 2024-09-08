package main

import (
	"log"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/config"
	"github.com/VanGoghDev/practicum-metrics/internal/agent/logger"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := make([]*analysis.Analyzer, 0, 5)

	// добавим стандартные проверки из analysis
	mychecks = append(mychecks, printf.Analyzer, shadow.Analyzer, structtag.Analyzer, errcheck.Analyzer)

	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name == "S1001" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)

	cfg, err := config.Load(&config.EnvReader{}, &config.FlagReader{})
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// logger
	zlog, err := logger.New(cfg.Loglevel)
	if err != nil {
		log.Fatal("failed to init logger %w", err)
	}

	zlog.Info("Logger init")

	agent := app.New(zlog, cfg)
	err = agent.RunApp()
	if err != nil {
		log.Fatalf("failed to run an app: %v", err)
	}
}
