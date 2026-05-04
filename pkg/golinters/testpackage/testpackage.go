package testpackage

import (
	"strings"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/maratori/testpackage/pkg/testpackage"
)

func New(settings *config.TestpackageSettings) *goanalysis.Linter {
	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			testpackage.SkipRegexpFlagName:    settings.SkipRegexp,
			testpackage.AllowPackagesFlagName: strings.Join(settings.AllowPackages, ","),
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(testpackage.NewAnalyzer()).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
