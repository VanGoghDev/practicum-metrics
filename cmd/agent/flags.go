package main

import (
	"flag"
)

var flagConsumerAddr string
var flagReportInterval int64
var flagPollInterval int64

func parseFlags() {
	flag.StringVar(&flagConsumerAddr, "a", "localhost:8080", "address and port where to send metrics")
	flag.Int64Var(&flagReportInterval, "r", 10, "report interval (interval of requests to consumer, in seconds)")
	flag.Int64Var(&flagPollInterval, "p", 2, "poll interval (interval of metrics fetch, in seconds)")
	flag.Parse()
}
