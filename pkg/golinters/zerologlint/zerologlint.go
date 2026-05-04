package zerologlint

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/ykadowak/zerologlint"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(zerologlint.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
