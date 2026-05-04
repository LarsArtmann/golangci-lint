package recvcheck

import (
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/raeperd/recvcheck"
)

func New(settings *config.RecvcheckSettings) *goanalysis.Linter {
	var cfg recvcheck.Settings

	if settings != nil {
		cfg.DisableBuiltin = settings.DisableBuiltin
		cfg.Exclusions = settings.Exclusions
	}

	return goanalysis.
		NewLinterFromAnalyzer(recvcheck.NewAnalyzer(cfg)).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
