package fatcontext

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"go.augendre.info/fatcontext/pkg/analyzer"
)

func New(settings *config.FatcontextSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			analyzer.FlagCheckStructPointers: settings.CheckStructPointers,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewAnalyzer()).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
