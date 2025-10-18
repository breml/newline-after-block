# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`newline-after-block` is a Go static analysis tool (linter) that enforces blank lines after block statements to improve code
readability. It's built using the `golang.org/x/tools/go/analysis` framework and can be used standalone or integrated with
golangci-lint. The analyzer supports automatic fixes that can insert missing blank lines via the `-fix` flag or through IDE
quick-fix suggestions.

## Architecture

The project has a simple but well-organized structure:

- **`newline-after-block.go`**: Core analyzer implementation
  - Defines the `Analyzer` using the `analysis.Analyzer` framework
  - `run()` function inspects AST nodes looking for `BlockStmt`, `SwitchStmt`, `TypeSwitchStmt`, and `SelectStmt` nodes
  - `checkStatements()` validates statement sequences for proper blank line spacing
  - `checkCaseClauses()` validates spacing between case clauses in switch/select statements
  - `needsNewlineAfter()` determines which statement types require blank lines (if without else, for, range, switch, type switch, select, defer)
  - `getBlockEnd()` extracts the end position of block statement bodies
  - `createDiagnosticWithFix()` creates diagnostics with suggested fixes to automatically insert blank lines
  - `findEndOfLine()` determines the correct position to insert newlines (handles inline comments)
  - `isErrorCheckIfStmt()` detects the `if err != nil` pattern for defer exceptions
  - `isErrNotNilPattern()` helper for error pattern matching
  - `isDeferStmt()` identifies defer statements

- **`cmd/newline-after-block/main.go`**: Command-line entry point
  - Uses `singlechecker.Main()` to create a standalone linter binary
  - Minimal wrapper around the analyzer

- **Test structure**: Uses `analysistest` framework
  - Test cases are in `testdata/src/` organized by package name
  - `testdata/src/blockstatements/` - tests for block statements (if, for, switch, etc.)
  - `testdata/src/caseclauses/` - tests for case clause spacing within switch/select statements
  - `testdata/src/structliterals/` - tests ensuring composite literals are not flagged
  - `testdata/src/deferpattern/` - tests for defer statement patterns after error checks
  - Tests use special `// want "..."` comments to verify expected diagnostics
  - Golden files (`.go.golden`) contain expected output after applying automatic fixes
  - `analysistest.RunWithSuggestedFixes()` verifies fixes produce correct output

## Key Linting Rules

The analyzer enforces blank lines after these block statements:

- `if` statements (only when not followed by `else`)
- `for` loops and `range` loops
- `switch` and type `switch` statements
- `select` statements

Additionally, the analyzer enforces blank lines between case clauses:

- Each case block within `switch`, type `switch`, and `select` statements must be followed by a blank line
- Exception: The last case block does not require a blank line before the closing brace
- Empty case blocks are skipped

It correctly ignores:

- Blocks at the end of statement lists (implicit)
- `if` statements followed by `else` or `else if`
- Composite literals (struct, array, slice, map literals)

Special handling for `defer` statements:

- `defer` statements can immediately follow error-checking `if err != nil` blocks without blank lines (idiomatic Go cleanup pattern)
- Multiple consecutive `defer` statements do not require blank lines between them
- A blank line IS required after `defer` statement(s) before any non-defer statement

## Autofix Capability

The analyzer provides automatic fix suggestions that can insert the required blank lines:

- **Command-line**: Use the `-fix` flag with the standalone binary to apply fixes automatically
- **IDE integration**: Editors with gopls support (e.g., VSCode) show "Quick Fix" suggestions
- **Implementation**: Each diagnostic includes a `SuggestedFix` with a `TextEdit` that inserts a newline at the correct position
- **Inline comments**: The fix correctly handles inline comments by inserting the newline after them
- **Idempotent**: Fixes can be applied multiple times without adverse effects
- **Best practice**: Apply fixes only to code that's already been formatted with `gofmt` or `gofumpt`

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
task lint-markdown            # Run markdownlint-cli2 on all markdown files
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

## GitHub Integration

For all GitHub-related tasks, use the `gh` command-line tool:

- **Pull Requests**: Use `gh pr create`, `gh pr view`, `gh pr list`, etc.
- **Issues**: Use `gh issue create`, `gh issue view`, `gh issue list`, etc.
- **Repository Information**: Use `gh repo view` to get repository details
- **GitHub URLs**: When given a GitHub URL, use the appropriate `gh` command to retrieve the information

Examples:

```bash
gh pr create --title "feat: Add new feature" --body "Description"
gh pr view 123
gh issue list --state open
gh repo view
```

The `gh` CLI provides a consistent and reliable interface for GitHub operations and should be preferred over direct API calls or web scraping.

## Commit Messages and PR Descriptions

### Commit Messages

Commit messages should be brief and to the point, following the Conventional Commits specification. Use the following prefixes:

- `fix` - Bug fixes
- `feat` - New features
- `docs` - Documentation changes
- `chore` - Maintenance tasks, dependency updates, etc.
- `pep` - Project Enhancement Proposals (used when adding new proposals to `admin/pep`)

Format: `<type>: <brief description>`

Examples:

```text
fix: Correct case clause spacing detection
feat: Add support for nested switch statements
docs: Update installation instructions
chore: Update dependencies to latest versions
pep: Add proposal for multi-file analysis support
```

### PR Descriptions

Pull request descriptions should be brief and to the point. Focus on:

- What changes were made (bullet points preferred)
- Why the changes were necessary (if not obvious)
- Avoid lengthy explanations unless the changes are complex

## Code Update Requirements

Whenever the code base is updated, ensure all of the following requirements are met:

1. **Tests**: All new code must have corresponding tests that verify the new functionality
2. **Formatting**: Code must be formatted using `gofumpt`
3. **Linting**: Code must be fully compliant with the golangci-lint configuration
4. **Verification**: Run `task verify` - it must complete without errors
5. **Documentation**:
   - Update `README.md` if the changes affect user-facing functionality or usage
   - Review and update `CLAUDE.md` if the changes affect architecture, development workflow, or coding standards

### Temporary Files

For temporary files or test artifacts, use the `.scratch/tmp/` directory.

## Adding Test Cases

When adding new test cases, create or modify files in `testdata/src/<packagename>/`. Use `analysistest` comment directives:

- `// want "message"` on the line where a diagnostic is expected
- The test framework will verify that the analyzer reports the expected diagnostic at that location
- For autofix testing, create a corresponding `.go.golden` file with the expected output after applying the fix

Example test pattern:

```go
if condition { // want "missing newline after block statement"
    doSomething()
}
nextStatement()
```

Corresponding `.go.golden` file:

```go
if condition {
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
