package inamedparam

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/macabu/inamedparam"
)

func New(settings *config.INamedParamSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"skip-single-param": settings.SkipSingleParam,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(inamedparam.Analyzer).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
