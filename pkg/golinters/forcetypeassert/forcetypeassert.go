package forcetypeassert

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/gostaticanalysis/forcetypeassert"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(forcetypeassert.Analyzer).
		WithDesc("Find forced type assertions").
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
