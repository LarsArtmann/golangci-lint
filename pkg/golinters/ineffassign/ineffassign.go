package ineffassign

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
)

func New(settings *config.IneffassignSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"check-escaping-errors": settings.CheckEscapingErrors,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(ineffassign.Analyzer).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
