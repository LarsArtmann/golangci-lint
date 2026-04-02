# AGENTS.md — golangci-lint

## Project Overview

**golangci-lint** is a fast Go linters runner that aggregates 60+ Go linters, runs them in parallel, and provides a unified configuration and output format. It is built on top of Go's `go/analysis` framework with a custom runner that adds caching, lazy-loading, and memory optimizations.

- **Module**: `github.com/golangci/golangci-lint/v2` (Go module path uses `/v2`)
- **Go version**: `go 1.25.0` minimum (always latest-1 Go version)
- **License**: GPL-3.0
- **Entry point**: `cmd/golangci-lint/main.go`
- **CLI framework**: `spf13/cobra`

## Essential Commands

### Build

```bash
make build                  # Build golangci-lint binary
go build -o golangci-lint ./cmd/golangci-lint   # Direct build
```

### Test

```bash
make test                   # Build + run golangci-lint on itself + go test ./... (requires CGO_ENABLED=1)
GL_TEST_RUN=1 go test -v -parallel 2 ./...    # Run tests directly
```

**Integration tests for a specific linter**:
```bash
# Test a specific linter's integration tests by testdata file name:
make test_integration T=bodyclose.go

# Test a specific fix test:
make test_integration_fix T=multiple-issues-fix.go
```

### Lint (self-lint)

```bash
make build
GL_TEST_RUN=1 ./golangci-lint run -v           # Run golangci-lint on itself
```

The project's own config is `.golangci.yml` at the repo root.

### Run a single linter manually on testdata

```bash
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=bodyclose ./pkg/golinters/bodyclose/testdata/bodyclose.go
```

### Generate / Documentation

```bash
make fast_generate                            # Generate github-action-config.json
make website_expand_templates                 # Expand docs templates
make docs_serve                               # Serve docs site locally
```

### Dependencies

```bash
go mod tidy && go mod verify                  # Tidy and verify modules
```

CI checks that `go.mod` and `go.sum` are clean after `go mod tidy`.

## Code Organization

```
cmd/golangci-lint/           # Main entry point (main.go, plugins.go)
pkg/commands/                # Cobra command definitions (run, linters, formatters, cache, config, migrate, etc.)
pkg/config/                  # Configuration structs, loading, and validation
pkg/goanalysis/              # Custom go/analysis runner (the core execution engine)
pkg/golinters/               # Linter wrappers — one subdirectory per linter
  <lintername>/              # Each contains:
    <lintername>.go          #   Linter adapter (returns *goanalysis.Linter)
    <lintername>_integration_test.go  # Integration test entry point
    testdata/                # Test fixture files (.go source + optional .yml config)
pkg/goformatters/            # Formatter wrappers (gofmt, goimports, gofumpt, gci, golines, etc.)
pkg/lint/                    # Core linting orchestration
  lintersdb/                 #   Linter database: builder, manager, validator
  linter/                    #   Linter/Config types, Context
pkg/result/                  # Issue types and result processors
  processors/                #   Post-processing pipeline (nolint, exclusion rules, severity, etc.)
pkg/printers/                # Output formatters (json, text, checkstyle, sarif, junit, etc.)
pkg/logutils/                # Logging utilities
pkg/fsutils/                 # Filesystem utilities
pkg/exitcodes/               # Exit code constants
pkg/goutil/                  # Go version utilities
internal/                    # Internal packages (cache, go utilities, x/tools vendored code)
test/                        # End-to-end integration tests
  testshared/                # Shared test utilities (runner builder, analysis, directives parser)
  testdata/                  # E2E test fixtures
scripts/                     # Build/release scripts, website generators
jsonschema/                  # JSON Schema files for configuration validation
docs/                        # Hugo-based documentation site
third_party/                 # Third-party licenses
```

## Architecture (Execution Flow)

1. **Init** — Config loaded from file/flags via `config.Loader`. Linter database built by `lintersdb.LinterBuilder`.
2. **Load Packages** — `go/packages` loads file info, types, and AST based on the union of all enabled linters' `LoadMode`.
3. **Run Linters** — Go/analysis-based linters are merged into a `MetaLinter` for parallel execution and caching. The `unused` linter is excluded from merging due to memory.
4. **Postprocess Issues** — Processors filter/transform issues: nolint directives, exclusion rules, severity, path shortening, deduplication.
5. **Print Issues** — Issues are output via the configured printer format.

## Adding a New Linter

This is one of the most common tasks. Follow these steps:

### 1. Create the linter adapter

Create `pkg/golinters/<lintername>/<lintername>.go`:

```go
package <lintername>

import (
    "github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
    return goanalysis.
        NewLinterFromAnalyzer(analyzer).
        WithLoadMode(goanalysis.LoadModeTypesInfo)
}
```

If the linter needs configuration, accept a settings struct:
```go
func New(settings *config.SomeSettings) *goanalysis.Linter { ... }
```

For linters with custom run logic (not just wrapping an analyzer), use `WithContextSetter` and `WithIssuesReporter`:
```go
return goanalysis.
    NewLinterFromAnalyzer(analyzer).
    WithContextSetter(func(lintCtx *linter.Context) {
        analyzer.Run = func(pass *analysis.Pass) (any, error) { ... }
    }).
    WithIssuesReporter(func(*linter.Context) []*goanalysis.Issue {
        return resIssues
    }).
    WithLoadMode(goanalysis.LoadModeTypesInfo)
```

### 2. Register in builder_linter.go

Add to `pkg/lint/lintersdb/builder_linter.go` in the `Build` method, **alphabetically (case-insensitive)**:

```go
linter.NewConfig(<lintername>.New(&cfg.Linters.Settings.SomeSettings)).
    WithSince("v2.X.0").       // Next minor version
    WithLoadForGoAnalysis().    // Only if LoadModeTypesInfo is needed
    WithURL("https://github.com/..."),
```

- Use `WithLoadForGoAnalysis()` if the linter uses `LoadModeTypesInfo`.
- Do NOT use `WithLoadForGoAnalysis()` if the linter uses `LoadModeSyntax` or `LoadModeNone`.
- Import the linter package alphabetically in the import block.

### 3. Add configuration (if needed)

- Add a settings struct to `pkg/config/linters_settings.go` with `mapstructure` tags.
- Add the field to the `LintersSettings` struct.
- Add defaults to `defaultLintersSettings`.
- Update `.golangci.next.reference.yml` (NOT `.golangci.reference.yml`).

### 4. Add tests

- Create `pkg/golinters/<lintername>/<lintername>_integration_test.go`:
  ```go
  func TestFromTestdata(t *testing.T) {
      integration.RunTestdata(t)
  }
  ```
- Create test data file `pkg/golinters/<lintername>/testdata/<lintername>.go` with test directives.
- For tests with custom config, also add a `.yml` file in testdata/.

### 5. Test directives format

Test files use `//golangcitest:` directives at the top:

```go
//golangcitest:args -Ebodyclose
//golangcitest:config_path testdata/whitespace.yml
//golangcitest:expected_linter whitespace
//golangcitest:expected_exitcode 1
package testdata

func Example() {
    // ... code that triggers the linter ...
    resp, _ := http.Get("https://example.com") // want "response body must be closed"
}
```

- `//golangcitest:args` — CLI arguments (e.g., `-Ebodyclose` to enable a linter)
- `//golangcitest:config_path` — Path to a config YAML file (relative to testdata dir)
- `//golangcitest:expected_linter` — Which linter should produce the diagnostic
- `//golangcitest:expected_exitcode` — Expected exit code (default: 1 = `exitcodes.IssuesFound`)

Expected diagnostics use `// want "pattern"` comments (inspired by `go/analysis/analysistest`).

### 6. Update `.golangci.next.reference.yml`

Add the linter configuration (alphabetical order). Do NOT edit `.golangci.reference.yml`.

## Key Patterns and Conventions

### Linter implementation patterns

There are three common patterns for linter adapters:

1. **Simple analyzer wrapper** — The linter just wraps an existing `analysis.Analyzer`:
   ```go
   func New() *goanalysis.Linter {
       return goanalysis.NewLinterFromAnalyzer(someanalyzer.Analyzer).
           WithLoadMode(goanalysis.LoadModeTypesInfo)
   }
   ```

2. **With configuration** — Linter accepts config settings and uses `WithContextSetter`:
   ```go
   func New(settings *config.SomeSettings) *goanalysis.Linter {
       var mu sync.Mutex
       var resIssues []*goanalysis.Issue
       analyzer := &analysis.Analyzer{ Name: "name", Doc: "...", Run: goanalysis.DummyRun }
       return goanalysis.NewLinterFromAnalyzer(analyzer).
           WithContextSetter(func(lintCtx *linter.Context) {
               analyzer.Run = func(pass *analysis.Pass) (any, error) { ... }
           }).
           WithIssuesReporter(func(*linter.Context) []*goanalysis.Issue { return resIssues }).
           WithLoadMode(goanalysis.LoadModeTypesInfo)
   }
   ```

3. **Simple no-config, no-types** — For AST-only linters:
   ```go
   func New() *goanalysis.Linter {
       return goanalysis.NewLinterFromAnalyzer(analyzer).
           WithLoadMode(goanalysis.LoadModeSyntax)
   }
   ```

### Go/analysis Load Modes

- `LoadModeNone` — No package data needed (e.g., `typecheck` is special)
- `LoadModeSyntax` — AST only, no type information. Fast.
- `LoadModeTypesInfo` — Full type information. Requires `WithLoadForGoAnalysis()` in builder. Marked as slow linter.
- `LoadModeWholeProgram` — Whole program analysis (e.g., `unused`). Not merged into MetaLinter.

### Configuration loading

- Config structs use `mapstructure` tags for YAML deserialization.
- Config is loaded by `pkg/config/loader.go` from `.golangci.yml` files and CLI flags.
- Placeholder support: `${envVAR}`, `~` home dir expansion.

### Logging

- Use `pkg/logutils` for logging. The project has its own logging abstraction (`logutils.Log`).
- Logging is only allowed through `logutils.Log`, NOT through `logrus` directly (enforced by depguard).
- Debug tags: `logutils.Debug(key)` returns a debug function for specific subsystems.

### Error handling

- Standard `errors` package and `fmt.Errorf` with `%w` wrapping.
- Do NOT use `github.com/pkg/errors` (enforced by depguard).

### Imports / formatting

- `interface{}` is automatically rewritten to `any` (enforced by gofmt config).
- `goimports` with local prefix `github.com/golangci/golangci-lint/v2`.
- Import groups: stdlib, then third-party, then project-internal.

## Testing

### Test environment variables

- `GL_TEST_RUN=1` — Must be set when running tests (enables verbose linter status output, used by test infrastructure).
- `GL_KEEP_TEMP_FILES=1` — Keep temporary test config files for debugging.

### Integration test flow

Integration tests work by:
1. Building the `golangci-lint` binary.
2. Running it against testdata `.go` files with specific directives.
3. Parsing `// want "pattern"` comments as expected diagnostics.
4. Comparing actual JSON output against expected patterns.

### Running specific tests

```bash
# All tests
make test

# Specific integration test
make test_integration T=bodyclose.go

# Specific Go test
GL_TEST_RUN=1 go test -v -run TestFromTestdata ./pkg/golinters/bodyclose/

# End-to-end tests
GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdata

# Specific E2E testdata file
GL_TEST_RUN=1 go test -v ./test -count 1 -run TestSourcesFromTestdata/bodyclose.go
```

### Test file naming

- `<lintername>.go` — Default test (no special config)
- `<lintername>_cgo.go` — CGO variant (same test, different build tag)
- `<lintername>_custom.go` + `<lintername>_custom.yml` — Tests with custom configuration
- `<lintername>_fix.go` — Fix mode tests (in `testdata/fix/in/` directory)

## Important Gotchas

1. **CGO_ENABLED=1 required for tests** — The Makefile sets `CGO_ENABLED=1` for the `test` target. Some testdata files use CGO build tags.

2. **`GL_TEST_RUN=1` is required** — Tests check this environment variable. Without it, some test behaviors differ.

3. **Module path is `/v2`** — All imports use `github.com/golangci/golangci-lint/v2/...`.

4. **Linters must be alphabetical** — In `builder_linter.go`, imports and linter registrations must be sorted alphabetically (case-insensitive).

5. **Only edit `.golangci.next.reference.yml`** — Never edit `.golangci.reference.yml` (that's the current released version's reference).

6. **Never edit the project's own `.golangci.yml`** for new linters — The project's self-lint config should not be modified when adding a new linter (per the new-linter-checklist).

7. **`nolintlint` runs last** — It must be the last linter in the processing pipeline (enforced by `linter.LastLinter` constant).

8. **The `unused` linter is special** — It uses `LoadModeWholeProgram` and is NOT merged into the MetaLinter due to high memory usage.

9. **Parallel test execution** — Integration tests use `t.Parallel()`. The test runner handles binary building and temp file cleanup.

10. **Build tags in test files** — Testdata files can use `//go:build` directives. The test infrastructure evaluates them and skips files that don't match the current platform/Go version.

11. **`go.mod` and `go.sum` must stay clean** — CI runs `go mod tidy` and checks for diffs.

12. **The Go version in `go.mod` is restrictive** — It's always set to the latest-1 Go version and should only be changed by maintainers when adding support for a new Go version.

## CI Workflows

- **`pr-tests.yml`** — Runs `make test` on Windows, macOS, and Linux (Ubuntu amd64 + arm64) with Go 1.25 and 1.26.
- **`pr-checks.yml`** — Checks `go mod tidy` is clean, tests CLI subcommands, and checks the install script.
- **`golangci-lint` action** — The repo also tests itself with the `golangci/golangci-lint-action` GitHub Action.
