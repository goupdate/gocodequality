package checks

import (
	"codequality"
	"go/ast"
)

const (
	// DefaultMaxComplexity is the default maximum cyclomatic complexity
	DefaultMaxComplexity = 10
)

type cyclomaticComplexityCheck struct {
	maxComplexity int
}

// CyclomaticComplexity creates a check that flags functions with complexity > 10
func CyclomaticComplexity() codequality.Check {
	return &cyclomaticComplexityCheck{maxComplexity: DefaultMaxComplexity}
}

func (c *cyclomaticComplexityCheck) Name() string {
	return "cyclomatic-complexity"
}

func (c *cyclomaticComplexityCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}

			complexity := 1 // Base complexity
			c.countComplexity(fd.Body, &complexity)

			if complexity > c.maxComplexity {
				pos := f.FSet.Position(fd.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Function '" + fd.Name.Name + "' has cyclomatic complexity " + itoa(complexity) + " (max " + itoa(c.maxComplexity) + ")",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}
		}
	}

	return issues
}

func (c *cyclomaticComplexityCheck) countComplexity(node ast.Node, complexity *int) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.CaseClause, *ast.SelectStmt, *ast.BinaryExpr:
			*complexity++
		}
		return true
	})
}
