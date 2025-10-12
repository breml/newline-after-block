# newline-after-block

A Go linter that enforces blank lines after block statements to improve code readability.

## Overview

`newline-after-block` is a static analysis tool that checks for the presence of blank lines after block statements in Go code. This helps maintain consistent code formatting and improves readability by visually separating logical blocks of code.

## Features

- Detects missing newlines after block statements (`if`, `for`, `switch`, `select`, etc.)
- Ignores composite literals (struct, array, slice, and map literals)
- Skips checks for blocks at the end of functions
- Respects `else` and `else if` clauses
- Provides clear, actionable error messages with file and line number references

## Installation

### Using `go install`

```bash
go install github.com/breml/newline-after-block/cmd/newline-after-block@latest
```

### Building from source

```bash
git clone https://github.com/breml/newline-after-block.git
cd newline-after-block
go build -o bin/newline-after-block ./cmd/newline-after-block
```

or if [`task`](https://taskfile.dev/) is installed:

```bash
task build
```

## Usage

### Command Line

Run the linter on your Go packages:

```bash
newline-after-block ./...
```

Run on specific files or directories:

```bash
newline-after-block ./pkg/mypackage
newline-after-block ./cmd/myapp/main.go
```

### Integration with golangci-lint

For integration with [golangci-lint](https://golangci-lint.run/), follow the instructions in
[Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/) and add the following to your `.golangci.yml`:

```yaml
linters:
  enable:
    - newline-after-block
  custom:
    newline-after-block:
      path: /path/to/newline-after-block
      description: Checks for newline after block statements
      original-url: https://github.com/breml/newline-after-block
```

## Rules

### Requires newline after

The linter enforces a blank line after these block statements:

- `if` statements (when not followed by `else`)
- `for` loops
- `range` loops
- `switch` statements
- `type switch` statements
- `select` statements

### Does NOT require newline after

The linter does not enforce newlines in these cases:

- Blocks at the end of functions
- `if` statements followed by `else` or `else if`
- Blocks followed by closing braces (e.g., end of another block)
- Composite literals (struct, array, slice, map literals)

## Examples

### Bad (will trigger linter)

```go
func example() {
    if condition {
        doSomething()
    } // missing blank line after if block
    nextStatement()
}

func loop() {
    for i := 0; i < 10; i++ {
        process(i)
    } // missing blank line after for block
    fmt.Println("done")
}
```

### Good (passes linter)

```go
func example() {
    if condition {
        doSomething()
    }

    nextStatement()  // Blank line present
}

func loop() {
    for i := 0; i < 10; i++ {
        process(i)
    }

    fmt.Println("done")  // Blank line present
}

func endOfFunction() {
    if condition {
        doSomething()
    }
    // No blank line needed - end of function
}

func elseClause() {
    if condition {
        doSomething()
    } else {
        doSomethingElse()
    }
    // No blank line needed after else
}

func structLiteral() {
    p := Person{
        Name: "John",
        Age:  30,
    }
    fmt.Println(p)  // No blank line needed - composite literal
}
```

## Development

### Prerequisites

- Go 1.25 or later
- [Task](https://taskfile.dev/) (optional, for running tasks)
- [golangci-lint](https://golangci-lint.run/) (for linting)

### Available Tasks

```bash
task build          # Build the linter binary
task test           # Run all tests
task test-verbose   # Run tests with verbose output
task test-coverage  # Run tests with coverage report
task lint           # Run golangci-lint
task lint-fix       # Run golangci-lint with auto-fix
task run            # Run the linter on the project itself
task clean          # Clean build artifacts
task verify         # Run all verification tasks (build, test, lint)
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Or using task
task test
```

### Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `task verify`
6. Submit a pull request

## Background

This linter was created to replace a shell script that used `sed` and `grep` to check for newlines after closing braces. The shell script had limitations:

- Could not distinguish between block statements and composite literals
- Did not handle ending of nested blocks correctly

The new Go-based linter uses the Go AST (Abstract Syntax Tree) to accurately identify block statements and apply the rules correctly.

## Author

Copyright 2025 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[MIT License](LICENSE)
