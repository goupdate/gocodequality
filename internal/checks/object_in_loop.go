package checks

import (
	"codequality"
	"go/ast"
)

type objectInLoopCheck struct{}

// ObjectInLoop creates a new objectinloop check
func ObjectInLoop() codequality.Check {
	return &objectInLoopCheck{}
}

func (c *objectInLoopCheck) Name() string {
	return "object-in-loop"
}

func (c *objectInLoopCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			forStmt, ok := n.(*ast.ForStmt)
			if !ok {
				return true
			}

			// Check for make(), new(), or time.Now() inside loop
			ast.Inspect(forStmt.Body, func(inner ast.Node) bool {
				call, ok := inner.(*ast.CallExpr)
				if !ok {
					return true
				}

				funcName := getCallName(call)
				if funcName == "make" || funcName == "new" || funcName == "time.Now" {
					pos := f.FSet.Position(call.Pos())
					issues = append(issues, codequality.Issue{
						Rule:     c.Name(),
						Severity: codequality.SeverityWarning,
						Message:  "Allocation '" + funcName + "()' inside loop - consider moving outside",
						File:     f.Path,
						Line:     pos.Line,
						Column:   pos.Column,
					})
				}

				return true
			})

			return true
		})
	}

	return issues
}

func getCallName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		if x, ok := fn.X.(*ast.Ident); ok {
			return x.Name + "." + fn.Sel.Name
		}
		return fn.Sel.Name
	}
	return ""
}
