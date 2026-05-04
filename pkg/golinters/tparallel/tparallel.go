package tparallel

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/moricho/tparallel"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(tparallel.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
