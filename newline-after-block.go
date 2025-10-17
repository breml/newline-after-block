// Package newlineafterblock provides a linter that checks for newlines after block statements.
package newlineafterblock

import (
	"go/ast"
	"go/token"
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

This rule also applies when a block statement is followed by a comment:
there should be a blank line between the block and the comment.

Additionally, this linter enforces blank lines between case clauses within
switch and select statements. Each case block (except the last) should be
followed by a blank line to improve readability.

Composite literals (struct/array/slice literals) and struct type definitions
are not considered block statements.`

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
		pass.Reportf(blockEnd, "missing newline after block statement")
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
			pass.Reportf(blockEnd, "missing newline after block statement")
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
			pass.Reportf(blockEnd, "missing newline after block statement")
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
		pass.Reportf(lastStmtEnd, "missing newline after case block")
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
			pass.Reportf(endPos, "missing newline after case block")
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
		pass.Reportf(lastStmtEnd, "missing newline after case block")
	}
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
	}

	return false
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
	}

	return token.NoPos
}
