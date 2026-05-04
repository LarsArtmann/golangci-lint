package noctx

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/sonatard/noctx"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(noctx.Analyzer).
		WithDesc("Detects function and method with missing usage of context.Context").
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
