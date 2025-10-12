# PEP 0004: Exclude File Patterns

Currently all files are included when linting for newlines after blocks. In Go projects, it is common to have generated files, that should not be checked for this rule.
This PEP proposes to add a configuration option to the analyzer to exclude certain file patterns from being checked.

The configuration option will be a list of regex patterns that will be matched against the relative file paths. If a file path matches any of the patterns, it will be excluded from the analysis.
If an invalid regex pattern is provided, the analyzer will return an error and exit with a non-zero status code.

The cli tool will be updated to accept a new flag `--exclude` (short flag `-e`) that can be provided multiple times to form a list of regex patterns. Use the `flag` package from the Go standard library to parse the flags. Implement a custom type that satisfies the `flag.Value` interface to handle this.

Update the testdata with a new file that matches one of the exclude patterns to ensure it is not checked. Add a test file @testdata/src/blockstatements/blockstatements_excluded.go, which contains block statements with following statements without newlines and ensure it is excluded from the analysis.
