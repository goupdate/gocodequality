package checks

import (
	"codequality"
	"go/ast"
	"go/token"
	"strings"
)

type duplicateImportsCheck struct{}

// DuplicateImports creates a new duplicateimports check
func DuplicateImports() codequality.Check {
	return &duplicateImportsCheck{}
}

func (c *duplicateImportsCheck) Name() string {
	return "duplicate-imports"
}

func (c *duplicateImportsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		seen := make(map[string]token.Pos)

		for _, decl := range f.AST.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.IMPORT {
				continue
			}

			for _, spec := range gd.Specs {
				is, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}

				path := strings.Trim(is.Path.Value, `"`)

				if prevPos, exists := seen[path]; exists {
					pos := f.FSet.Position(is.Pos())
					issues = append(issues, codequality.Issue{
						Rule:     c.Name(),
						Severity: codequality.SeverityWarning,
						Message:  "Duplicate import '" + path + "' (first imported at line " + itoa(f.FSet.Position(prevPos).Line) + ")",
						File:     f.Path,
						Line:     pos.Line,
						Column:   pos.Column,
					})
				} else {
					seen[path] = is.Pos()
				}
			}
		}
	}

	return issues
}
