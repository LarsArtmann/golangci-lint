package exptostd

import (
	"testing"

	"github.com/golangci/golangci-lint/v2/test/testshared/integration"
	// Without this dependency, the analyzer tests related to fix fails.
	// The packages `slices` have been randomly chosen to import `golang.org/x/exp`.
	_ "golang.org/x/exp/slices"
)

func TestFromTestdata(t *testing.T) {
	integration.RunTestdata(t)
}

func TestFix(t *testing.T) {
	integration.RunFix(t)
}

func TestFixPathPrefix(t *testing.T) {
	integration.RunFixPathPrefix(t)
}
