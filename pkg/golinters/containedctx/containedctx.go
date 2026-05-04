package containedctx

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/sivchari/containedctx"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(containedctx.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
