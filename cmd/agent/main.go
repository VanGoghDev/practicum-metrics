package main

import (
	"time"

	"github.com/VanGoghDev/practicum-metrics/internal/agent/app"
)

func main() {
	parseFlags()
	app := app.New(flagConsumerAddr, time.Duration(flagReportInterval)*time.Second, time.Duration(flagPollInterval)*time.Second)
	app.MustRun()
}
