package checks

import (
	"codequality"
	"go/ast"
)

type panicUsageCheck struct{}

// PanicUsage creates a new panicusage check
func PanicUsage() codequality.Check {
	return &panicUsageCheck{}
}

func (c *panicUsageCheck) Name() string {
	return "panic-usage"
}

func (c *panicUsageCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if it's a panic() call
			if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				pos := f.FSet.Position(call.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Use of panic() detected - consider using proper error handling",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}

			return true
		})
	}

	return issues
}
