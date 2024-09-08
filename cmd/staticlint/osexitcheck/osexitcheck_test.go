package osexitcheck_test

import (
	"testing"

	"github.com/VanGoghDev/practicum-metrics/cmd/staticlint/osexitcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), osexitcheck.Analyzer, "./...")
}
