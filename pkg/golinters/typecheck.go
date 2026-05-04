package golinters

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewTypecheck() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(&analysis.Analyzer{
			Name: "typecheck",
			Doc:  "Like the front-end of a Go compiler, parses and type-checks Go code",
			Run:  goanalysis.DummyRun,
		}).
		WithLoadMode(goanalysis.LoadModeNone)
}
