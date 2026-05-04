package funcorder

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/manuelarte/funcorder/analyzer"
)

func New(settings *config.FuncOrderSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			analyzer.ConstructorCheckName:  settings.Constructor,
			analyzer.StructMethodCheckName: settings.StructMethod,
			analyzer.AlphabeticalCheckName: settings.Alphabetical,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewAnalyzer()).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
