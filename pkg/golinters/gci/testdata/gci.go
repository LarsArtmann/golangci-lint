//golangcitest:config_path testdata/gci.yml
package testdata

import (
	"errors"
	"fmt"

	gcicfg "github.com/daixiang0/gci/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"golang.org/x/tools/go/analysis" // want "File is not properly formatted"
)

func GoimportsLocalTest() {
	fmt.Print(errors.New("x"))
	_ = config.Config{}
	_ = analysis.Analyzer{}
	_ = gcicfg.BoolConfig{}
}
