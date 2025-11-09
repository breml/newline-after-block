// Package newlineafterblock provides a linter that checks for newlines after block statements.
package newlineafterblock

import (
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
)

// doc describes what the analyzer does.
const doc = `check for newline after block statements

This linter ensures that block statements (if, for, switch, select, etc.)
are followed by a blank line, unless:
- The block is at the end of a function
- The block is followed by an else/else if
- The block is followed by a closing brace
- The block is followed by another case/default in a switch/select
- The block is an error-checking if statement (if <error> != nil) followed by a defer

This rule also applies when a block statement is followed by a comment:
there should be a blank line between the block and the comment.

Additionally, this linter enforces blank lines between case clauses within
switch and select statements. Each case block (except the last) should be
followed by a blank line to improve readability.

Special handling for defer statements:
- Defer statements can immediately follow error-checking if statements (if <error> != nil)
  without a blank line (idiomatic Go pattern for cleanup)
- Error detection is type-based: any variable implementing the error interface is recognized
- Multiple consecutive defer statements do not require blank lines between them
- A blank line is required after defer statement(s) before any non-defer statement

Composite literals (struct/array/slice literals) and struct type definitions
are not considered block statements.

The analyzer provides automatic fix suggestions that insert the required blank
lines.`

type newlineafterblock struct {
	exclude excludePatterns
}

// New creates and returns a new newline-after-block analyzer instance.
func New() *analysis.Analyzer {
	nlab := newlineafterblock{}

	analyzer := &analysis.Analyzer{
		Name: "newlineafterblock",
		Doc:  doc,
		Run:  nlab.run,
	}

	// Register flags on this analyzer instance.
	analyzer.Flags.Var(&nlab.exclude, "exclude", "regex pattern to exclude files from analysis")
	analyzer.Flags.Var(&nlab.exclude, "e", "regex pattern to exclude files from analysis (shorthand)")

	return analyzer
}

func (n *newlineafterblock) run(pass *analysis.Pass) (any, error) {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}

	for _, file := range pass.Files {
		if n.shouldSkipFile(pass, file, wd) {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			n.inspectNode(pass, file, node)
			return true
		})
	}

	return nil, nil
}

// shouldSkipFile determines if a file should be skipped based on exclude patterns.
func (n *newlineafterblock) shouldSkipFile(pass *analysis.Pass, file *ast.File, wd string) bool {
	relPath, err := filepath.Rel(wd, pass.Fset.Position(file.Package).Filename)
	if err != nil {
		relPath = pass.Fset.Position(file.Package).Filename
	}

	return n.exclude.matches(relPath)
}

// inspectNode inspects an AST node and performs appropriate checks.
func (n *newlineafterblock) inspectNode(pass *analysis.Pass, file *ast.File, node ast.Node) {
	switch n := node.(type) {
	case *ast.BlockStmt:
		checkStatements(pass, file, n.List)

	case *ast.CaseClause:
		checkStatements(pass, file, n.Body)

	case *ast.SwitchStmt:
		if n.Body != nil {
			checkCaseClauses(pass, file, n.Body.List)
		}

	case *ast.TypeSwitchStmt:
		if n.Body != nil {
			checkCaseClauses(pass, file, n.Body.List)
		}

	case *ast.SelectStmt:
		if n.Body != nil {
			checkCommClauses(pass, file, n.Body.List)
		}
	}
}

// checkStatements checks a sequence of statements for missing newlines after blocks.
func checkStatements(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	for i := 0; i < len(stmts)-1; i++ {
		checkStatementPair(pass, astFile, stmts[i], stmts[i+1])
	}

	// Also check the last statement if it's followed by a comment.
	if len(stmts) > 0 {
		checkLastStatement(pass, astFile, stmts[len(stmts)-1])
	}
}

// checkStatementPair checks if there's proper spacing between two consecutive statements.
func checkStatementPair(pass *analysis.Pass, astFile *ast.File, current, next ast.Stmt) {
	// Exception: Allow defer immediately after error-checking if statement.
	if isErrorCheckIfStmt(pass, current) && isDeferStmt(next) {
		return
	}

	// Exception: Allow consecutive defer statements without blank line.
	if isDeferStmt(current) && isDeferStmt(next) {
		return
	}

	if !needsNewlineAfter(current) {
		return
	}

	blockEnd := getBlockEnd(current)
	if blockEnd == token.NoPos {
		return
	}

	file := pass.Fset.File(blockEnd)
	if file == nil {
		return
	}

	blockEndLine := file.Line(blockEnd)
	nextLine := file.Line(next.Pos())

	// Check if there's a comment between the block and the next statement.
	foundComment := checkCommentBetween(pass, astFile, file, blockEnd, blockEndLine, next.Pos())

	// If no comment was found between the block and next statement,
	// check if the next statement is immediately after (no blank line).
	if !foundComment && nextLine == blockEndLine+1 {
		pass.Report(createDiagnosticWithFix(pass, blockEnd, "missing newline after block statement"))
	}
}

// checkCommentBetween checks for comments between a block end and the next statement.
// Returns true if a non-inline comment was found.
func checkCommentBetween(pass *analysis.Pass, astFile *ast.File, file *token.File, blockEnd token.Pos, blockEndLine int, nextPos token.Pos) bool {
	for _, commentGroup := range astFile.Comments {
		if commentGroup.Pos() <= blockEnd || commentGroup.Pos() >= nextPos {
			continue
		}

		commentLine := file.Line(commentGroup.Pos())
		// Skip inline comments (on the same line as the closing brace).
		if commentLine == blockEndLine {
			continue
		}

		// Found a comment on a different line.
		// If comment is on the next line (no blank line).
		if commentLine == blockEndLine+1 {
			pass.Report(createDiagnosticWithFix(pass, blockEnd, "missing newline after block statement"))
		}

		// Only check the first non-inline comment.
		return true
	}

	return false
}

// checkLastStatement checks if the last statement has proper spacing before any trailing comments.
func checkLastStatement(pass *analysis.Pass, astFile *ast.File, lastStmt ast.Stmt) {
	if !needsNewlineAfter(lastStmt) {
		return
	}

	blockEnd := getBlockEnd(lastStmt)
	if blockEnd == token.NoPos {
		return
	}

	file := pass.Fset.File(blockEnd)
	if file == nil {
		return
	}

	blockEndLine := file.Line(blockEnd)

	// Check if there's a comment after the last statement.
	checkTrailingComment(pass, astFile, file, blockEnd, blockEndLine)
}

// checkTrailingComment checks for comments after a block statement.
func checkTrailingComment(pass *analysis.Pass, astFile *ast.File, file *token.File, blockEnd token.Pos, blockEndLine int) {
	for _, commentGroup := range astFile.Comments {
		if commentGroup.Pos() <= blockEnd {
			continue
		}

		commentLine := file.Line(commentGroup.Pos())
		// Skip inline comments (on the same line as the closing brace).
		if commentLine == blockEndLine {
			continue
		}

		// If comment is on the next line (no blank line).
		if commentLine == blockEndLine+1 {
			pass.Report(createDiagnosticWithFix(pass, blockEnd, "missing newline after block statement"))
		}

		// Only check the first comment after the block.
		break
	}
}

// checkCaseClauses checks that case clauses in switch/select statements are properly spaced.
// Each case clause (except the last) should be followed by a blank line.
func checkCaseClauses(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	caseClauses := extractCaseClauses(stmts)
	if len(caseClauses) < 2 {
		return
	}

	// Check spacing between consecutive case clauses.
	for i := 0; i < len(caseClauses)-1; i++ {
		checkCaseClauseSpacing(pass, astFile, caseClauses[i], caseClauses[i+1])
	}
}

// extractCaseClauses filters statements to only CaseClause nodes.
func extractCaseClauses(stmts []ast.Stmt) []*ast.CaseClause {
	var caseClauses []*ast.CaseClause
	for _, stmt := range stmts {
		if caseClause, ok := stmt.(*ast.CaseClause); ok {
			caseClauses = append(caseClauses, caseClause)
		}
	}

	return caseClauses
}

// checkCaseClauseSpacing checks spacing between two consecutive case clauses.
func checkCaseClauseSpacing(pass *analysis.Pass, astFile *ast.File, current, next *ast.CaseClause) {
	// Skip empty case clauses (no body statements).
	if len(current.Body) == 0 {
		return
	}

	lastStmt := current.Body[len(current.Body)-1]
	lastStmtEnd := lastStmt.End()

	file := pass.Fset.File(lastStmtEnd)
	if file == nil {
		return
	}

	lastStmtLine := file.Line(lastStmtEnd)
	nextCaseLine := file.Line(next.Pos())

	// Check if there's a comment between the last statement and the next case.
	foundComment := checkClauseComment(pass, astFile, file, lastStmtEnd, lastStmtLine, next.Pos())

	// If no comment was found, check if the next case is immediately after.
	if !foundComment && nextCaseLine == lastStmtLine+1 {
		pass.Report(createDiagnosticWithFix(pass, lastStmtEnd, "missing newline after case block"))
	}
}

// checkClauseComment checks for comments between two clause positions and reports violations.
// Returns true if a non-inline comment was found.
func checkClauseComment(pass *analysis.Pass, astFile *ast.File, file *token.File, endPos token.Pos, endLine int, nextPos token.Pos) bool {
	for _, commentGroup := range astFile.Comments {
		commentPos := commentGroup.Pos()
		if commentPos <= endPos || commentPos >= nextPos {
			continue
		}

		commentLine := file.Line(commentPos)
		// Skip inline comments (on the same line as the end position).
		if commentLine == endLine {
			continue
		}

		// If comment is on the next line (no blank line).
		if commentLine == endLine+1 {
			pass.Report(createDiagnosticWithFix(pass, endPos, "missing newline after case block"))
		}

		// Only check the first non-inline comment.
		return true
	}

	return false
}

// checkCommClauses checks that comm clauses in select statements are properly spaced.
// Each comm clause (except the last) should be followed by a blank line.
// CommClause is used for select statements, similar to CaseClause for switch statements.
func checkCommClauses(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	commClauses := extractCommClauses(stmts)
	if len(commClauses) < 2 {
		return
	}

	// Check spacing between consecutive comm clauses.
	for i := 0; i < len(commClauses)-1; i++ {
		checkCommClauseSpacing(pass, astFile, commClauses[i], commClauses[i+1])
	}
}

// extractCommClauses filters statements to only CommClause nodes.
func extractCommClauses(stmts []ast.Stmt) []*ast.CommClause {
	var commClauses []*ast.CommClause
	for _, stmt := range stmts {
		if commClause, ok := stmt.(*ast.CommClause); ok {
			commClauses = append(commClauses, commClause)
		}
	}

	return commClauses
}

// checkCommClauseSpacing checks spacing between two consecutive comm clauses.
func checkCommClauseSpacing(pass *analysis.Pass, astFile *ast.File, current, next *ast.CommClause) {
	// Skip empty comm clauses (no body statements).
	if len(current.Body) == 0 {
		return
	}

	lastStmt := current.Body[len(current.Body)-1]
	lastStmtEnd := lastStmt.End()

	file := pass.Fset.File(lastStmtEnd)
	if file == nil {
		return
	}

	lastStmtLine := file.Line(lastStmtEnd)
	nextCommLine := file.Line(next.Pos())

	// Check if there's a comment between the last statement and the next comm.
	foundComment := checkClauseComment(pass, astFile, file, lastStmtEnd, lastStmtLine, next.Pos())

	// If no comment was found, check if the next comm is immediately after.
	if !foundComment && nextCommLine == lastStmtLine+1 {
		pass.Report(createDiagnosticWithFix(pass, lastStmtEnd, "missing newline after case block"))
	}
}

// checkAssignStmt checks if an assignment statement contains a function literal.
func checkAssignStmt(s *ast.AssignStmt) *ast.FuncLit {
	for _, expr := range s.Rhs {
		if funcLit := extractFuncLit(expr); funcLit != nil {
			return funcLit
		}
	}

	return nil
}

// checkDeclStmt checks if a declaration statement contains a function literal.
func checkDeclStmt(s *ast.DeclStmt) *ast.FuncLit {
	genDecl, ok := s.Decl.(*ast.GenDecl)
	if !ok {
		return nil
	}

	for _, spec := range genDecl.Specs {
		if funcLit := checkValueSpec(spec); funcLit != nil {
			return funcLit
		}
	}

	return nil
}

// checkValueSpec checks if a value spec contains a function literal.
func checkValueSpec(spec ast.Spec) *ast.FuncLit {
	valueSpec, ok := spec.(*ast.ValueSpec)
	if !ok {
		return nil
	}

	for _, value := range valueSpec.Values {
		if funcLit := extractFuncLit(value); funcLit != nil {
			return funcLit
		}
	}

	return nil
}

// extractFuncLit extracts a function literal from an expression.
// Only returns function literals that are NOT immediately invoked,
// since invoked function literals end with ), not }.
func extractFuncLit(expr ast.Expr) *ast.FuncLit {
	switch e := expr.(type) {
	case *ast.FuncLit:
		// Direct function literal: func() {}
		return e

	case *ast.CallExpr:
		// Don't flag immediately invoked function literals: func() {}()
		// These end with ), not }, so they shouldn't require a newline
		return nil
	}

	return nil
}

// needsNewlineAfter determines if a statement needs a newline after it.
func needsNewlineAfter(stmt ast.Stmt) bool {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		// If statement with else: check the else branch.
		// If statement without else: needs newline
		if s.Else != nil {
			// The else branch itself needs a newline after it.
			return true
		}

		return true

	case *ast.ForStmt, *ast.RangeStmt:
		return true

	case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true

	case *ast.AssignStmt:
		return checkAssignStmt(s) != nil

	case *ast.DeclStmt:
		return checkDeclStmt(s) != nil

	case *ast.DeferStmt:
		// Defer statements need newlines when followed by non-defer statements.
		// The exception (consecutive defers) is handled in checkStatementPair.
		return true
	}

	return false
}

// isErrorCheckIfStmt checks if an if statement matches the pattern "if <error> != nil".
func isErrorCheckIfStmt(pass *analysis.Pass, stmt ast.Stmt) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return false
	}

	// Check if the condition is a binary expression.
	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		return false
	}

	// Check if the operator is !=.
	if binaryExpr.Op != token.NEQ {
		return false
	}

	// Check if one operand is a variable implementing error interface and the other is nil.
	return isErrNotNilPattern(pass, binaryExpr.X, binaryExpr.Y) || isErrNotNilPattern(pass, binaryExpr.Y, binaryExpr.X)
}

// isErrNotNilPattern checks if x is a variable implementing the error interface and y is nil.
func isErrNotNilPattern(pass *analysis.Pass, x, y ast.Expr) bool {
	ident, ok := x.(*ast.Ident)
	if !ok {
		return false
	}

	// Check if y is nil.
	nilIdent, ok := y.(*ast.Ident)
	if !ok || nilIdent.Name != "nil" {
		return false
	}

	// Check if x has a type that implements the error interface.
	if pass.TypesInfo == nil {
		return false
	}

	typ := pass.TypesInfo.TypeOf(ident)
	if typ == nil {
		return false
	}

	// Check if the type implements the error interface.
	return implementsError(typ)
}

// implementsError checks if a type implements the error interface using types.Implements.
func implementsError(typ types.Type) bool {
	errorObj := types.Universe.Lookup("error")
	if errorObj == nil {
		return false
	}

	errorObjType := errorObj.Type()
	if errorObjType == nil {
		return false
	}

	underlying := errorObjType.Underlying()
	if underlying == nil {
		return false
	}

	errorType, ok := underlying.(*types.Interface)
	if !ok {
		return false
	}

	// Check both value and pointer receiver cases.
	return types.Implements(typ, errorType) || types.Implements(types.NewPointer(typ), errorType)
}

// isDeferStmt checks if a statement is a defer statement.
func isDeferStmt(stmt ast.Stmt) bool {
	_, ok := stmt.(*ast.DeferStmt)
	return ok
}

// getBlockEnd returns the end position of a block statement's body.
func getBlockEnd(stmt ast.Stmt) token.Pos {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		// If there's an else branch, return the end of the entire if-else chain.
		if s.Else != nil {
			return getBlockEnd(s.Else)
		}

		// Otherwise return the end of the if body.
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.BlockStmt:
		// Handle else blocks (which are BlockStmt nodes).
		return s.End()

	case *ast.ForStmt:
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.RangeStmt:
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.SwitchStmt:
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.TypeSwitchStmt:
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.SelectStmt:
		if s.Body != nil {
			return s.Body.End()
		}

	case *ast.AssignStmt:
		if funcLit := checkAssignStmt(s); funcLit != nil && funcLit.Body != nil {
			return funcLit.Body.End()
		}

	case *ast.DeclStmt:
		if funcLit := checkDeclStmt(s); funcLit != nil && funcLit.Body != nil {
			return funcLit.Body.End()
		}

	case *ast.DeferStmt:
		// For defer statements, return the end position of the statement.
		return s.End()
	}

	return token.NoPos
}

// findEndOfLine returns the position at the end of the line containing pos.
// This handles inline comments automatically since we insert at end of current line.
func findEndOfLine(file *token.File, pos token.Pos) token.Pos {
	line := file.Line(pos)

	// If not the last line, return the start of the next line
	// (which is right after the newline character of the current line).
	if line < file.LineCount() {
		return file.LineStart(line + 1)
	}

	// Last line: return end of file.
	return token.Pos(file.Base() + file.Size())
}

// createDiagnosticWithFix creates a diagnostic with a suggested fix to insert a blank line.
func createDiagnosticWithFix(pass *analysis.Pass, blockEnd token.Pos, message string) analysis.Diagnostic {
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
