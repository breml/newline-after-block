# PEP: Add Autofix Capability to newline-after-block

## Overview

This proposal outlines the implementation plan for adding automatic fix (autofix) capability to the `newline-after-block` linter. The autofix feature will allow the linter to automatically insert missing newlines after block statements via command-line flag.

## Motivation

Currently, `newline-after-block` only reports violations. Users must manually add the required blank lines. Adding autofix capability will:

1. **Improve developer experience**: Developers can automatically fix all violations with a single command
2. **Facilitate adoption**: Easier to apply to large existing codebases
3. **Align with modern tooling**: Most contemporary linters support autofix (e.g., `gofmt`, `goimports`)

## Technical Approach

### 1. Core Implementation Changes

#### 1.1 Switch from `Reportf` to `Report` with Diagnostics

Currently, violations are reported using:

```go
pass.Reportf(blockEnd, "missing newline after block statement")
```

This needs to be changed to:

```go
pass.Report(analysis.Diagnostic{
    Pos:     blockEnd,
    Message: "missing newline after block statement",
    SuggestedFixes: []analysis.SuggestedFix{
        {
            Message: "Insert blank line",
            TextEdits: []analysis.TextEdit{
                {
                    Pos:     insertPos,
                    End:     insertPos,
                    NewText: []byte("\n"),
                },
            },
        },
    },
})
```

#### 1.2 Determine Correct Insertion Position

The fix must insert a newline at the correct position. Key considerations:

- **Insert position**: After the closing brace of the block statement (at `blockEnd`)
- **Account for inline comments**: If there's a comment on the same line as the closing brace, insert after the comment
- **Preserve existing formatting**: Don't disrupt indentation or other whitespace

**Algorithm**:

1. Get the position of the closing brace (`blockEnd`)
2. Scan forward to the end of the line to handle inline comments
3. Insert a newline character at the end of the line

**Implementation sketch**:

```go
func findEndOfLine(file *token.File, pos token.Pos) token.Pos {
    line := file.Line(pos)

    // Not the last line in the file.
    // Return the position just before the next line's start (the newline character position).
    // This handles inline comments automatically since we insert at end of current line.
    if line < file.LineCount() {
        return file.LineStart(line + 1) - 1
    }

    // Last line in the file.
    return token.Pos(file.Base() + file.Size())
}
```

#### 1.3 Modify Report Call Sites

Replace all four instances of `pass.Reportf()` with `pass.Report()`:

1. **checkStatementPair**: Line 122 - when next statement is immediately after block
2. **checkCommentBetween**: Line 143 - when comment is immediately after block
3. **checkTrailingComment**: Line 190 - when trailing comment is immediately after block (last statement case)

Each call site needs:

- Create `analysis.Diagnostic` struct
- Calculate correct insertion position
- Create `analysis.TextEdit` with newline insertion
- Wrap in `analysis.SuggestedFix`

### 2. Helper Functions

Create helper functions to keep code DRY and maintainable:

#### 2.1 `createDiagnosticWithFix`

Since the current implementation only uses a fixed message string (`"missing newline after block statement"`), a simple `createDiagnosticWithFix` function is sufficient. However, if future enhancements require formatted messages, this could be extended to `createDiagnosticWithFixf` accepting `fmt.Sprintf`-style formatting:

```go
func createDiagnosticWithFix(
    pass *analysis.Pass,
    blockEnd token.Pos,
    message string,
) analysis.Diagnostic {
    file := pass.Fset.File(blockEnd)
    if file == nil {
        // Fallback: return diagnostic without fix
        return analysis.Diagnostic{
            Pos:     blockEnd,
            Message: message,
        }
    }

    // Find the end of the line containing blockEnd
    insertPos := findEndOfLine(file, blockEnd)

    return analysis.Diagnostic{
        Pos:     blockEnd,
        Message: message,
        SuggestedFixes: []analysis.SuggestedFix{
            {
                Message: "Insert blank line after block statement",
                TextEdits: []analysis.TextEdit{
                    {
                        Pos:     insertPos,
                        End:     insertPos,
                        NewText: []byte("\n"),
                    },
                },
            },
        },
    }
}
```

**Note**: For the current implementation, the simple version without formatting is recommended. A formatted variant (`createDiagnosticWithFixf`) is not needed unless different error messages are required for different violation types.

#### 2.2 `findEndOfLine`

```go
func findEndOfLine(file *token.File, pos token.Pos) token.Pos {
    line := file.Line(pos)

    // If not the last line, return position just before next line's start
    if line < file.LineCount() {
        return file.LineStart(line + 1) - 1
    }

    // Last line: return end of file
    return token.Pos(file.Base() + file.Size())
}
```

### 3. Testing Strategy

#### 3.1 Update Test Framework

The Go `analysistest` package supports testing suggested fixes using `analysistest.RunWithSuggestedFixes()`.

**Test file naming convention**:

- Input files: `testdata/src/package/file.go`
- Expected output: `testdata/src/package/file.go.golden`

#### 3.2 Create Golden Files

For each test case that has a `// want "..."` comment, create a corresponding `.golden` file showing the expected result after applying the fix.

**Example**:

`testdata/src/blockstatements/blockstatements.go`:

```go
func ifStatementWithoutNewline() {
    x := 5
    if x > 0 {
        fmt.Println("positive")
    } // want "missing newline after block statement"
    fmt.Println("next statement")
}
```

`testdata/src/blockstatements/blockstatements.go.golden`:

```go
func ifStatementWithoutNewline() {
    x := 5
    if x > 0 {
        fmt.Println("positive")
    }

    fmt.Println("next statement")
}
```

#### 3.3 Add Autofix Tests

Add a new test function:

```go
func TestAnalyzerWithFixes(t *testing.T) {
    analyzer := newlineafterblock.New()
    testdata := analysistest.TestData()
    analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "blockstatements")
}
```

This test will:

1. Run the analyzer on test files
2. Apply all suggested fixes
3. Compare the result with `.golden` files
4. Fail if results don't match

### 4. Edge Cases and Considerations

#### 4.1 Inline Comments

When a closing brace has an inline comment:

```go
if x > 0 {
    fmt.Println("positive")
} // some comment
fmt.Println("next")
```

The newline should be inserted after the comment:

```go
if x > 0 {
    fmt.Println("positive")
} // some comment

fmt.Println("next")
```

This is already handled by inserting at end of line.

#### 4.2 Independent and Idempotent Fixes

All fixes are designed to be independent and idempotent. Since this linter targets well-formatted Go code (users are expected to use `gofmt` or `gofumpt`), edge cases like multiple block statements on the same line are not a concern. Each fix inserts exactly one blank line, and applying the same fix multiple times has no additional effect.

#### 4.3 Platform Line Endings

The implementation should use `\n` for consistency with Go conventions. The `go fmt` tool normalizes line endings, so this should not cause issues.

#### 4.4 Preserving Subsequent Blank Lines

If there are already multiple blank lines after a block (not a violation), the fix shouldn't affect them. Since we're only flagging cases with zero blank lines, this isn't an issue.

### 5. Implementation Steps

#### Phase 1: Core Implementation

1. Add `createDiagnosticWithFix` helper function
2. Add `findEndOfLine` helper function
3. Update `checkStatementPair` to use `pass.Report()` with suggested fix
4. Update `checkCommentBetween` to use `pass.Report()` with suggested fix
5. Update `checkTrailingComment` to use `pass.Report()` with suggested fix

#### Phase 2: Testing

1. Create golden files for all existing test cases in `testdata/src/blockstatements/`
2. Create golden files for test cases in `testdata/src/comments/`
3. Add `TestAnalyzerWithFixes` for blockstatements package
4. Add `TestAnalyzerWithFixes` for comments package
5. Verify all tests pass

#### Phase 3: Documentation

1. Update `README.md` to document autofix capability
2. Update doc string in analyzer to mention autofix support
3. Add usage examples for standalone binary with `-fix` flag
4. Document that fixes should only be applied to already formatted code (via `gofmt` or `gofumpt`)

#### Phase 4: Validation

**Note**: Validation is performed externally and manually, and is therefore out of scope for this proposal. The following activities are recommended for post-implementation validation:

1. Run the analyzer with fixes on this codebase itself
2. Performance testing on large codebases

### 6. Compatibility Considerations

- **Backward compatibility**: The change is fully backward compatible. Users who don't use autofix features will see no change in behavior.
- **Command-line flag**: A `-fix` flag will be added to the standalone binary to apply suggested fixes automatically.

### 7. Success Criteria

The autofix implementation is successful when:

1. All existing tests continue to pass
2. New tests with golden files verify correct fix application
3. Running the analyzer with `-fix` flag on already properly formatted code (via `gofmt` or `gofumpt`) produces valid, formatted Go code with the required blank lines inserted
4. VSCode with gopls shows "Quick Fix" suggestions
5. Performance impact is negligible (< 5% slower than current implementation)

## References

- [Go Analysis Package Documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Analysis Testing Package](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest)
- [SuggestedFix Examples](https://github.com/golang/tools/blob/master/go/analysis/passes/printf/printf.go)
