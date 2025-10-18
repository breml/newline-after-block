# PEP 0009: Enforce Newline After Multi-line Function Signature

## Overview

This proposal outlines extending the `newline-after-block` linter to enforce a blank line at the beginning of a function body when the function signature spans multiple lines. This improves code readability by providing visual separation between the function declaration and its implementation.

## Motivation

When function signatures are long and span multiple lines (due to multiple parameters, type constraints, or complex return types), there's reduced visual separation between the signature and the function body. Adding a blank line after the opening brace provides:

1. **Improved readability**: Clear visual distinction between function declaration and implementation
2. **Consistency**: Aligns with the linter's philosophy of enforcing blank lines for better code structure

## Examples

### Incorrect (No blank line after multi-line signature)

```go
func processUserRequest(
    userID int,
    email string,
    options map[string]interface{},
) error {
    // Implementation starts immediately
    user := getUser(userID)
    if user == nil {
        return fmt.Errorf("user not found")
    }
    return nil
}
```

### Correct (Blank line inserted after opening brace)

```go
func processUserRequest(
    userID int,
    email string,
    options map[string]interface{},
) error {

    // Implementation starts with blank line
    user := getUser(userID)
    if user == nil {
        return fmt.Errorf("user not found")
    }
    return nil
}
```

### When This Rule Applies

The rule applies only when:
- The function signature spans **multiple lines** (i.e., the opening `{` is not on the same line as the `func` keyword)
- The function has a body (not an interface method without implementation)

### When This Rule Does NOT Apply

The rule does not apply to:
- Single-line function signatures (standard case)
- Empty function bodies
- Interface method declarations (they have no body)
- Anonymous functions (lambdas) that are on a single line

## Technical Approach

### 1. Detection Logic

To determine if a function has a multi-line signature:

1. Get the position of the `func` keyword
2. Get the position of the opening brace `{`
3. Compare their line numbers using the token file
4. If `{` is on a different line than `func`, the signature is multi-line
5. Check if there's already a blank line after the opening brace

### 2. Implementation Strategy

#### 2.1 New Checker Function

Add a new function to detect multi-line function signatures:

```go
// checkMultilineFunctionSignature checks if a function with a multi-line
// signature has a blank line after the opening brace.
func checkMultilineFunctionSignature(
    pass *analysis.Pass,
    funcDecl *ast.FuncDecl,
) {
    // Get function keyword position
    funcPos := funcDecl.Pos()

    // Get opening brace position
    bracePos := funcDecl.Body.Lbrace

    // Get file information
    file := pass.Fset.File(funcPos)
    if file == nil {
        return
    }

    // Compare line numbers
    funcLine := file.Line(funcPos)
    braceLine := file.Line(bracePos)

    // Multi-line signature: opening brace is on different line than func keyword
    if braceLine == funcLine {
        return // Single-line signature, skip
    }

    // Check if there's a blank line after the opening brace
    if hasBlankLineAfterBrace(pass, funcDecl.Body) {
        return // Already has blank line
    }

    // Report violation with suggested fix
    diag := createDiagnosticWithFix(
        pass,
        bracePos,
        "missing newline after multi-line function signature",
    )
    pass.Report(diag)
}
```

#### 2.2 Blank Line Detection

Add a helper function to check for existing blank line:

```go
// hasBlankLineAfterBrace checks if there's a blank line immediately after
// the opening brace of a block.
func hasBlankLineAfterBrace(pass *analysis.Pass, block *ast.BlockStmt) bool {
    if len(block.List) == 0 {
        return true // Empty block is considered valid
    }

    file := pass.Fset.File(block.Lbrace)
    if file == nil {
        return true
    }

    braceLine := file.Line(block.Lbrace)
    firstStmtLine := file.Line(block.List[0].Pos())

    // If first statement is more than 1 line away, there's a blank line
    return firstStmtLine > braceLine+1
}
```

#### 2.3 Integration into Main Run Function

Modify the `run()` function to check all function declarations:

```go
// In the run() function, add:
for _, decl := range f.Decls {
    if funcDecl, ok := decl.(*ast.FuncDecl); ok {
        if funcDecl.Body != nil {
            checkMultilineFunctionSignature(pass, funcDecl)
        }
    }
}
```

### 3. Testing Strategy

#### 3.1 New Test Cases

Create test file at `testdata/src/multilinefunction/multilinefunction.go`:

```go
package multilinefunction

// Single-line signature - should not trigger
func singleLine(x int) error {
    return nil
}

// Multi-line signature without blank line - should trigger
func multilineNoBlank(
    userID int,
    email string,
) error { // want "missing newline after multi-line function signature"
    return nil
}

// Multi-line signature with blank line - should pass
func multilineWithBlank(
    userID int,
    email string,
) error {

    return nil
}

// Method receiver on same line as func - still multi-line signature
func (r *Receiver) methodWithMultiline(
    param1 string,
    param2 int,
) error { // want "missing newline after multi-line function signature"
    return nil
}

// Empty function body
func emptyFunc(
    param string,
) {
}

// Named return values
func withNamedReturns(
    x int,
) (result string, err error) { // want "missing newline after multi-line function signature"
    return "", nil
}
```

#### 3.2 Golden File

Create `testdata/src/multilinefunction/multilinefunction.go.golden`:

```go
package multilinefunction

// Single-line signature - should not trigger
func singleLine(x int) error {
    return nil
}

// Multi-line signature without blank line - should trigger
func multilineNoBlank(
    userID int,
    email string,
) error {

    return nil
}

// Multi-line signature with blank line - should pass
func multilineWithBlank(
    userID int,
    email string,
) error {

    return nil
}

// Method receiver on same line as func - still multi-line signature
func (r *Receiver) methodWithMultiline(
    param1 string,
    param2 int,
) error {

    return nil
}

// Empty function body
func emptyFunc(
    param string,
) {
}

// Named return values
func withNamedReturns(
    x int,
) (result string, err error) {

    return "", nil
}
```

#### 3.3 Test Code

Add new test function to `newline_after_block_test.go`:

```go
func TestMultilineFunctionSignature(t *testing.T) {
    analyzer := newlineafterblock.New()
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, analyzer, "multilinefunction")
}

func TestMultilineFunctionSignatureWithFixes(t *testing.T) {
    analyzer := newlineafterblock.New()
    testdata := analysistest.TestData()
    analysistest.RunWithSuggestedFixes(t, testdata, analyzer, "multilinefunction")
}
```

### 4. Edge Cases and Considerations

#### 4.1 Empty Function Bodies

Empty function bodies (just `{ }` with no statements) should not trigger the violation, as there's nowhere to add a blank line.

#### 4.2 Interface Method Declarations

Interface method declarations don't have bodies and should be skipped entirely. This is already handled by checking `if funcDecl.Body != nil`.

#### 4.3 Receiver Types

For methods with receiver types, the rule still applies. The receiver is part of the signature, and we check based on the `func` keyword position:

```go
// This is multi-line even though receiver is on same line
func (r *Receiver) method(
    param string,
) error { // Multi-line, should have blank line after
    return nil
}
```

#### 4.4 Generic Type Constraints

For functions with type constraints (Go 1.18+), the signature is still based on the `func` keyword position:

```go
func genericFunc[T comparable](
    x T,
) T { // Multi-line, should have blank line after
    return x
}
```

#### 4.5 Formatting with gofmt

The Go formatter (`gofmt`) does not remove blank lines inside function bodies, so this change is safe for already-formatted code. The fix is idempotent.

### 5. Implementation Steps

#### Phase 1: Core Implementation

1. Add `checkMultilineFunctionSignature` function to analyze function declarations
2. Add `hasBlankLineAfterBrace` helper function
3. Integrate checks into the `run()` function
4. Ensure fix suggestion uses existing `createDiagnosticWithFix` helper
5. Handle edge cases (empty bodies, interface methods)

#### Phase 2: Testing

1. Create test file at `testdata/src/multilinefunction/multilinefunction.go`
2. Create corresponding golden file
3. Add test functions to verify detection
4. Add test functions to verify fixes
5. Run `task verify` to ensure all tests and linting pass

#### Phase 3: Documentation

1. Update `README.md` to document the new rule
2. Add description to analyzer documentation
3. Document which cases trigger the rule and which don't
4. Add examples to README showing the new behavior
5. Update `CLAUDE.md` if necessary

### 6. Compatibility Considerations

- **Backward compatibility**: This is a new rule and does not affect existing diagnostics
- **Configuration**: If the linter supports rule disabling in the future, this rule should be independently disableable
- **Go version**: No version-specific constraints; applies to all supported Go versions

### 7. Success Criteria

The implementation is successful when:

1. All existing tests continue to pass
2. New tests correctly identify multi-line function signatures without blank lines
3. New tests verify that single-line signatures are not flagged
4. Golden files verify correct fix application for all test cases
5. Running `task verify` completes without errors
6. Running the analyzer with `-fix` on this codebase produces valid formatted Go code
7. Documentation is updated with examples and explanations
8. No false positives on interface method declarations or other special cases

## References

- [Go Function Declarations](https://golang.org/ref/spec#Function_declarations)
- [Go Analyzer Package Documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Analyzing AST Nodes](https://golang.org/pkg/go/ast/)
