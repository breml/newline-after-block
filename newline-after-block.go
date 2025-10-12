// Package newlineafterblock provides a linter that checks for newlines after block statements.
package newlineafterblock

import (
	"go/ast"
	"go/token"
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
	for _, file := range pass.Files {
		pos := pass.Fset.Position(file.Pos())
		relPath, err := filepath.Rel(filepath.Dir(pass.Fset.Position(file.Package).Filename), pos.Filename)
		if err != nil {
			relPath = pos.Filename
		}

		// Check if this file matches any exclude pattern.
		if n.exclude.matches(relPath) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			// Check BlockStmt nodes to find statement sequences.
			if block, ok := n.(*ast.BlockStmt); ok {
				checkStatements(pass, file, block.List)
			}

			// Check CaseClause nodes (switch/select case bodies).
			if caseClause, ok := n.(*ast.CaseClause); ok {
				checkStatements(pass, file, caseClause.Body)
			}

			return true
		})
	}

	return nil, nil
}

// checkStatements checks a sequence of statements for missing newlines after blocks.
func checkStatements(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	for i := 0; i < len(stmts)-1; i++ {
		current := stmts[i]
		next := stmts[i+1]

		// Check if current statement is a block statement that needs a newline.
		if needsNewlineAfter(current) {
			// Get the position of the closing brace.
			blockEnd := getBlockEnd(current)
			if blockEnd == token.NoPos {
				continue
			}

			file := pass.Fset.File(blockEnd)
			if file == nil {
				continue
			}

			blockEndLine := file.Line(blockEnd)
			nextLine := file.Line(next.Pos())

			// Check if there's a comment between the block and the next statement.
			// We need to skip inline comments (comments on the same line as the closing brace).
			foundComment := false
			for _, commentGroup := range astFile.Comments {
				if commentGroup.Pos() <= blockEnd || commentGroup.Pos() >= next.Pos() {
					continue
				}

				commentLine := file.Line(commentGroup.Pos())
				// Skip inline comments (on the same line as the closing brace).
				if commentLine == blockEndLine {
					continue
				}

				// Found a comment on a different line.
				foundComment = true
				// If comment is on the next line (no blank line).
				if commentLine == blockEndLine+1 {
					pass.Reportf(blockEnd, "missing newline after block statement")
				}

				// Only check the first non-inline comment.
				break
			}

			// If no comment was found between the block and next statement,
			// check if the next statement is immediately after (no blank line).
			if !foundComment && nextLine == blockEndLine+1 {
				pass.Reportf(blockEnd, "missing newline after block statement")
			}
		}
	}

	// Also check the last statement if it's followed by a comment.
	if len(stmts) > 0 {
		lastStmt := stmts[len(stmts)-1]
		if needsNewlineAfter(lastStmt) {
			blockEnd := getBlockEnd(lastStmt)
			if blockEnd != token.NoPos {
				file := pass.Fset.File(blockEnd)
				if file != nil {
					blockEndLine := file.Line(blockEnd)
					// Check if there's a comment after the last statement.
					// Skip inline comments (on the same line as the closing brace).
					for _, commentGroup := range astFile.Comments {
						if commentGroup.Pos() > blockEnd {
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
				}
			}
		}
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
