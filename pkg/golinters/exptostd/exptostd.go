package exptostd

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/ldez/exptostd"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(exptostd.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
