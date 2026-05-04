package wastedassign

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/sanposhiho/wastedassign/v2"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(wastedassign.Analyzer).
		WithDesc("Finds wasted assignment statements").
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
