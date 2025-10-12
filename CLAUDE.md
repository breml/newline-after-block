# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`newline-after-block` is a Go static analysis tool (linter) that enforces blank lines after block statements to improve code readability. It's built using the `golang.org/x/tools/go/analysis` framework and can be used standalone or integrated with golangci-lint.

## Architecture

The project has a simple but well-organized structure:

- **`newline-after-block.go`**: Core analyzer implementation
  - Defines the `Analyzer` using the `analysis.Analyzer` framework
  - `run()` function inspects AST nodes looking for `BlockStmt` nodes
  - `checkStatements()` validates statement sequences for proper blank line spacing
  - `needsNewlineAfter()` determines which statement types require blank lines (if without else, for, range, switch, type switch, select)
  - `getBlockEnd()` extracts the end position of block statement bodies

- **`cmd/newline-after-block/main.go`**: Command-line entry point
  - Uses `singlechecker.Main()` to create a standalone linter binary
  - Minimal wrapper around the analyzer

- **Test structure**: Uses `analysistest` framework
  - Test cases are in `testdata/src/` organized by package name
  - `testdata/src/blockstatements/` - tests for block statements (if, for, switch, etc.)
  - `testdata/src/structliterals/` - tests ensuring composite literals are not flagged
  - Tests use special `// want "..."` comments to verify expected diagnostics

## Key Linting Rules

The analyzer enforces blank lines after these block statements:

- `if` statements (only when not followed by `else`)
- `for` loops and `range` loops
- `switch` and type `switch` statements
- `select` statements

It correctly ignores:

- Blocks at the end of statement lists (implicit)
- `if` statements followed by `else` or `else if`
- Composite literals (struct, array, slice, map literals)

## Development Commands

### Building

```bash
task build                    # Build to bin/newline-after-block
```

### Testing

```bash
task test                     # Run all tests
go test ./...                 # Direct test command
task test-coverage            # Generate coverage report (coverage.html)
```

### Linting

```bash
task lint                     # Run golangci-lint
task lint-fix                 # Run golangci-lint with auto-fix
```

### Running the linter

```bash
task run                      # Run on this project itself
go run ./cmd/newline-after-block ./...
```

### Full verification

```bash
task verify                   # Runs mod-tidy, build, test, lint
```

## Adding Test Cases

When adding new test cases, create or modify files in `testdata/src/<packagename>/`. Use `analysistest` comment directives:

- `// want "message"` on the line where a diagnostic is expected
- The test framework will verify that the analyzer reports the expected diagnostic at that location

Example test pattern:

```go
if condition { // want "missing newline after block statement"
    doSomething()
}
nextStatement()
```

## Coding Standards

When modifying code in this repository, adhere to the following standards:

- **Cognitive Complexity**: Functions must have a cognitive complexity of 20 or less (enforced by golangci-lint with gocognit)
  - To reduce complexity, extract nested logic into helper functions
  - Use early returns to reduce nesting
  - Break down complex functions into smaller, single-purpose functions
  - The refactoring approach used in this codebase demonstrates how to handle this constraint

## Go Version

Requires Go 1.25 or later (as specified in go.mod).
