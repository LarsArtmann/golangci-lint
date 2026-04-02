# Status Report: Fang v2 Integration for golangci-lint

**Date:** 2026-04-01 06:40  
**Branch:** `feat/fang-v2-integration`  
**PR:** https://github.com/golangci/golangci-lint/pull/6473  
**Issue:** https://github.com/golangci/golangci-lint/issues/6464  

---

## A) FULLY DONE

1. **Fang v2 integrated** — `fang.Execute()` replaces `c.cmd.Execute()` in `root.go`
2. **Custom brand theme** — `theme.go` with golangci-lint.run colors (blue `#2563eb`, teal `#0d9488`)
3. **`color.GreenString()` removed** from ~60 flag descriptions across 12 files
4. **Unused `color` imports removed** from 7 files (`flagsets.go`, `run.go`, `custom.go`, `fmt.go`, `migrate.go`, `help.go`, `config.go`)
5. **Duplicate error printing removed** from `main.go` (fang handles styled errors)
6. **Build passes** — `go build ./cmd/golangci-lint/` clean
7. **Tests pass** — `go test ./pkg/commands/...` all green
8. **`--help` styled** — All subcommands show themed help output
9. **Error output styled** — Unknown flags etc. show themed error box with "Try --help"
10. **`help linters` / `help formatters`** — Still work correctly with their `color.Green`/`color.Red` headers
11. **`version` subcommand** — Unchanged, `--json`/`--debug`/`--short` still work
12. **PR created and pushed** — 2 commits on branch, PR #6473 open

---

## B) PARTIALLY DONE (needs more work)

1. **`--version` output format changed** — Root `--version` now outputs `golangci-lint version (devel)` (fang's terse format) instead of the original rich format `golangci-lint has version (devel) built with go1.26.1 from (abc1234, modified: false, mod sum: "...") on 2026-03-28`. We removed `rootOptions.PrintVersion`, the `--version` flag from `rootCmd.Flags()`, and the `PrintVersion` handler from `RunE`. Fang now owns `--version` via `fang.WithVersion(c.info.Version)`.
2. **Theme comments** — `blue700` comment says `"primary hover"` which is a CSS/web concept, not CLI. Should be `"// error detail text"` or similar.

---

## C) NOT STARTED

1. **Visual testing** — Theme has not been tested on a real terminal (dark/light mode)
2. **`--color` flag interaction with fang** — `color.NoColor` from `fatih/color` may conflict with fang's own color detection
3. **Fang's hidden `man` command** — Fang adds a hidden `man` subcommand by default. Not investigated.
4. **Fang's completion generation** — Fang adds shell completions by default. Not investigated.
5. **CI testing** — Not verified in GitHub Actions pipeline
6. **Squash commits** — Branch has 2 commits, should be 1 clean commit
7. **Write tests** — No tests added for the new fang integration
8. **Update PR description** — Needs to reflect all issues found and fixed

---

## D) TOTALLY FUCKED UP

1. **`--version` regression** — The biggest issue. We replaced the custom `--version` handler (`PrintVersion`) with fang's `WithVersion()`. This changes the output from:
   ```
   golangci-lint has version 1.64.0 built with go1.26.1 from (abc1234, modified: false, mod sum: "...") on 2026-03-28
   ```
   to:
   ```
   golangci-lint version (devel)
   ```
   Users lose Go version, commit SHA, date, and build info. This breaks CI scripts that parse `--version` output for build metadata.

2. **We removed `rootOptions.PrintVersion`** and the `--version` flag from `rootCmd.Flags()` — this was unnecessary. Fang's `WithVersion()` sets `root.Version` which cobra uses for its built-in `--version` flag. But cobra's built-in `--version` just prints `<cmd> version <string>` — it doesn't call our `BuildInfo.String()`.

---

## E) WHAT WE SHOULD IMPROVE

1. **Restore `--version` to original behavior** — Use `fang.WithoutVersion()` and keep the custom `--version` flag + `PrintVersion` handler in `RunE`. Fang should only own help and error styling, not version output.
2. **Fix theme comments** — Replace CSS terminology (`primary hover`) with CLI-appropriate descriptions (`error detail text`).
3. **Squash into 1 clean commit** — Current 2 commits tell a messy story.
4. **Test on real terminal** — Verify theme looks good in both light and dark mode.
5. **Investigate `--color` flag conflict** — `color.NoColor` may interfere with fang's color detection.
6. **Investigate fang's hidden features** — `man` command, completions.
7. **Add tests** — At minimum, test `--version` output format to prevent future regressions.

---

## F) TOP 25 THINGS TO DO NEXT (priority order)

| # | Task | Impact | Effort |
|---|------|--------|--------|
| 1 | Restore `--version` to original rich output (`WithoutVersion()` + custom handler) | HIGH | small |
| 2 | Fix theme.go comments (remove CSS terminology) | LOW | tiny |
| 3 | Build, test, verify restored `--version` | HIGH | tiny |
| 4 | Commit the fix | HIGH | tiny |
| 5 | Squash all commits into 1 clean commit | MED | small |
| 6 | Force push | HIGH | tiny |
| 7 | Update PR description with all changes | MED | small |
| 8 | Run `go mod tidy` to fix lipgloss direct warning | LOW | tiny |
| 9 | Investigate `--color` flag conflict with fang | MED | med |
| 10 | Visual test theme on real terminal (dark + light) | MED | small |
| 11 | Investigate fang's hidden `man` command | LOW | small |
| 12 | Investigate fang's completion generation | LOW | small |
| 13 | Remove remaining `color.GreenString("auto-fix")` in `help_linters.go` | LOW | tiny |
| 14 | Write test for `--version` output format | MED | small |
| 15 | Run full CI pipeline | MED | med |
| 16 | Check if existing CI scripts parse `--version` output | MED | small |
| 17 | Consider adding `WithoutManpage()` to fang options | LOW | tiny |
| 18 | Consider adding `WithoutCompletions()` to fang options | LOW | tiny |
| 19 | Review PR for any other regressions | MED | small |
| 20 | Test `golangci-lint run` still works end-to-end | HIGH | small |
| 21 | Test piping output (`golangci-lint run | cat`) | MED | small |
| 22 | Verify no ANSI codes in piped/non-TTY output | MED | small |
| 23 | Check gopls diagnostics are clean | LOW | tiny |
| 24 | Final review of diff before merge | HIGH | small |
| 25 | Merge when approved | HIGH | tiny |

---

## G) TOP #1 QUESTION I CANNOT FIGURE OUT MYSELF

**Should `--version` on root use fang's terse format or keep the original rich format?**

- **Option A:** Keep `fang.WithVersion()` → terse `golangci-lint version X.Y.Z` (cleaner, fang-native)
- **Option B:** Use `fang.WithoutVersion()` + custom `--version` flag → rich `golangci-lint has version X.Y.Z built with go1.26.1 from (commit) on date` (preserves existing behavior)

My recommendation: **Option B** — use `fang.WithoutVersion()`. Fang is great for help/error styling, but the version output has specific build metadata (Go version, commit SHA, date) that fang's format doesn't support. The `version` subcommand with `--json`/`--debug`/`--short` already exists for structured access, but the root `--version` should keep its rich human-readable format for CI scripts and debugging.

This is a product decision though — I need your input.
