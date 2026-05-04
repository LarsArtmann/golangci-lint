package sqlclosecheck

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/ryanrolds/sqlclosecheck/pkg/analyzer"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewDeferOnlyAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
