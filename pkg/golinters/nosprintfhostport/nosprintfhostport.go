package nosprintfhostport

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/stbenjam/no-sprintf-host-port/pkg/analyzer"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.Analyzer).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
