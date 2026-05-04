package arangolint

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"go.augendre.info/arangolint/pkg/analyzer"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
