package checks

import (
	"codequality"
	"go/ast"
)

const (
	// DefaultMaxParams is the default maximum number of function parameters
	DefaultMaxParams = 5
)

type tooManyParamsCheck struct {
	maxParams int
}

// TooManyParams creates a check that flags functions with more than 5 parameters
func TooManyParams() codequality.Check {
	return &tooManyParamsCheck{maxParams: DefaultMaxParams}
}

func (c *tooManyParamsCheck) Name() string {
	return "too-many-params"
}

func (c *tooManyParamsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Type.Params == nil {
				continue
			}

			paramCount := 0
			for _, field := range fd.Type.Params.List {
				if len(field.Names) == 0 {
					paramCount++
				} else {
					paramCount += len(field.Names)
				}
			}

			if paramCount > c.maxParams {
				pos := f.FSet.Position(fd.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Function '" + fd.Name.Name + "' has " + itoa(paramCount) + " parameters (max " + itoa(c.maxParams) + ")",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}
		}
	}

	return issues
}
