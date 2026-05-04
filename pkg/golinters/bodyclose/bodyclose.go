package bodyclose

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/timakin/bodyclose/passes/bodyclose"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(bodyclose.Analyzer).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
