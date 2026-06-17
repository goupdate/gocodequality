package checks

import (
	"codequality"
	"go/ast"
)

type ignoredErrorsCheck struct{}

// IgnoredErrors creates a new ignorederrors check
func IgnoredErrors() codequality.Check {
	return &ignoredErrorsCheck{}
}

func (c *ignoredErrorsCheck) Name() string {
	return "ignored-errors"
}

func (c *ignoredErrorsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			assign, ok := n.(*ast.AssignStmt)
			if !ok {
				return true
			}

			// Check for _ = someFunc() pattern
			if len(assign.Lhs) == 1 {
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name == "_" {
					// Check if RHS is a function call
					if call, ok := assign.Rhs[0].(*ast.CallExpr); ok {
						pos := f.FSet.Position(assign.Pos())
						funcName := getFuncName(call)
						issues = append(issues, codequality.Issue{
							Rule:     c.Name(),
							Severity: codequality.SeverityWarning,
							Message:  "Error return value ignored for '" + funcName + "'",
							File:     f.Path,
							Line:     pos.Line,
							Column:   pos.Column,
						})
					}
				}
			}

			return true
		})
	}

	return issues
}

func getFuncName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		if x, ok := fn.X.(*ast.Ident); ok {
			return x.Name + "." + fn.Sel.Name
		}
		return fn.Sel.Name
	}
	return "unknown"
}
