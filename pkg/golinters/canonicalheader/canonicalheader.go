package canonicalheader

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/lasiar/canonicalheader"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(canonicalheader.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
