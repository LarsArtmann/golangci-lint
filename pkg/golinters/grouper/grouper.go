package grouper

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	grouper "github.com/leonklingele/grouper/pkg/analyzer"
)

func New(settings *config.GrouperSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"const-require-single-const":   settings.ConstRequireSingleConst,
			"const-require-grouping":       settings.ConstRequireGrouping,
			"import-require-single-import": settings.ImportRequireSingleImport,
			"import-require-grouping":      settings.ImportRequireGrouping,
			"type-require-single-type":     settings.TypeRequireSingleType,
			"type-require-grouping":        settings.TypeRequireGrouping,
			"var-require-single-var":       settings.VarRequireSingleVar,
			"var-require-grouping":         settings.VarRequireGrouping,
		}
	}
	return goanalysis.
		NewLinterFromAnalyzer(grouper.New()).
		WithDesc("Analyze expression groups.").
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
