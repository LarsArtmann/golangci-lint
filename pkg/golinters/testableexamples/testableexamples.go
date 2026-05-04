package testableexamples

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/maratori/testableexamples/pkg/testableexamples"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(testableexamples.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
