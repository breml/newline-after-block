# PEP 0001 -- Initial Version of newline-after-block linter

Create a linter that checks for a newline after block statements in Go code.

The linter is expected to replace the existing shell script used for this purpose, which simply checks for a newline after the final curly brace of block statements but does not handle composite literals correctly.

The current shell script is as follows:

```sh
#!/bin/sh -eu

echo "Checking that functional blocks are followed by newlines..."

# Check all .go files except the protobuf bindings (.pb.go)
files=$(git ls-files --cached --modified --others '*.go' ':!:*.pb.go' ':!:test/mini-oidc/storage/*.go' ':!:internal/server/network/*/schema/*/*.go')

exit_code=0
for file in $files; do
    # This oneliner has a few steps:
    # 1. sed:
    #     a. Check for lines that contain a single closing brace (plus whitespace).
    #     b. Move the pattern space window forward to the next line.
    #     c. Match lines that start with a word character. This allows for a closing brace on subsequent lines.
    #     d. Print the line number.
    # 2. xargs: Print the filename next to the line number of the matches (piped).
    # 3. If there were no matches, the file name without the line number is printed, use grep to filter it out.
    # 4. Replace the space with a colon to make a clickable link.
    RESULT=$(sed -n -e '/^\s*}\s*$/{n;/^\s*\w/{;=}}' "$file" | xargs -L 1 echo "$file" | grep -v '\.go$' | sed 's/ /:/g')
    if [ -n "${RESULT}" ]; then
        echo "${RESULT}"
        exit_code=1
    fi
done

exit $exit_code
```

## Specification

The linter should enforce the following rules:

1. There must be a blank line after every block statement (e.g., `if`, `for`, `switch`, `select`, etc.).
2. The linter should ignore newline requirements after composite literals (e.g., struct literals, array literals).
3. The linter should provide clear and actionable feedback to the user when a violation is detected.

## Implementation Details

The linter will be implemented as a standalone tool that can be integrated into existing Go development workflows. It will be built using the Go programming language and will leverage the **Go AST** (Abstract Syntax Tree) to analyze code structure.
The linter should follow the best practices outlined in the [New linters](https://golangci-lint.run/docs/contributing/new-linters/) guidelines from golangci-lint.
The linter should have decent test coverage to ensure reliability and maintainability.
The linter should be checked for linting issues using `golangci-lint`.
The repository should include a `Taskfile.yml` to facilitate common tasks such as building, testing, and running the linter.

## Repository structure

The repository for the linter will have the following structure:

```none
newline-after-block/
├── cmd/
│   └── newline-after-block/
│       └── main.go              # Entry point for the linter
├── newline-after-block.go       # Core linter logic
├── newline-after-block_test.go  # Tests for the linter
├── go.mod                       # Go module file
├── go.sum                       # Go module checksum file
├── Taskfile.yml                 # Taskfile for common tasks
├── testdata/                     # Directory for test data files
│   └── src/
│       ├── blockstatements/
│       └── structliterals/
├── .golangci.yml                # Configuration for golangci-lint
└── README.md                    # Documentation for the linter
```
