package internal

import (
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func FormatCode(code string) string {
	if strings.Contains(code, "`") {
		return code // TODO: properly escape or remove
	}

	return fmt.Sprintf("`%s`", code)
}

func GetGoFileNames(pass *analysis.Pass) []string {
	var filenames []string

	for _, f := range pass.Files {
		position, b := goanalysis.GetGoFilePosition(pass, f)
		if !b {
			continue
		}

		filenames = append(filenames, position.Filename)
	}

	return filenames
}
