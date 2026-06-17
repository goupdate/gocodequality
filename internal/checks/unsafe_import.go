package checks

import (
	"codequality"
	"go/ast"
)

type unsafeImportCheck struct{}

// UnsafeImport creates a new unsafeimport check
func UnsafeImport() codequality.Check {
	return &unsafeImportCheck{}
}

func (c *unsafeImportCheck) Name() string {
	return "unsafe-import"
}

func (c *unsafeImportCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, imp := range f.AST.Imports {
			if imp.Path.Value == `"unsafe"` {
				pos := f.FSet.Position(imp.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Import of 'unsafe' package detected - potential memory safety issues",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}
		}
	}

	return issues
}

// Ensure ast import is used
var _ = ast.NewIdent
