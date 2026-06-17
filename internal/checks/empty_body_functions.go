package checks

import (
	"codequality"
	"go/ast"
	"strings"
)

type emptyBodyFunctionsCheck struct{}

// EmptyBodyFunctions creates a new emptybodyfunctions check
func EmptyBodyFunctions() codequality.Check {
	return &emptyBodyFunctionsCheck{}
}

func (c *emptyBodyFunctionsCheck) Name() string {
	return "empty-body-functions"
}

func (c *emptyBodyFunctionsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			fd, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Skip interface methods (no body)
			if fd.Body == nil {
				return true
			}

			// Skip test/benchmark/example functions
			name := fd.Name.Name
			if strings.HasPrefix(name, "Test") ||
				strings.HasPrefix(name, "Benchmark") ||
				strings.HasPrefix(name, "Example") {
				return true
			}

			// Check if body is empty
			if len(fd.Body.List) == 0 {
				pos := f.FSet.Position(fd.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Function '" + name + "' has an empty body",
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
