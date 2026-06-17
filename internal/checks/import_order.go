package checks

import (
	"codequality"
	"go/ast"
	"go/token"
	"strings"
)

type importOrderCheck struct{}

// ImportOrder creates a new importorder check
func ImportOrder() codequality.Check {
	return &importOrderCheck{}
}

func (c *importOrderCheck) Name() string {
	return "import-order"
}

func (c *importOrderCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.IMPORT {
				continue
			}

			if len(gd.Specs) < 2 {
				continue
			}

			// Check if imports are grouped properly
			seenExternal := false
			for _, spec := range gd.Specs {
				is, ok := spec.(*ast.ImportSpec)
				if !ok {
					continue
				}

				path := strings.Trim(is.Path.Value, `"`)
				isStdlib := isStdlibPackage(path)

				if isStdlib && seenExternal {
					pos := f.FSet.Position(is.Pos())
					issues = append(issues, codequality.Issue{
						Rule:     c.Name(),
						Severity: codequality.SeverityInfo,
						Message:  "Standard library import '" + path + "' should be before external imports",
						File:     f.Path,
						Line:     pos.Line,
						Column:   pos.Column,
					})
				}

				if !isStdlib {
					seenExternal = true
				}
			}
		}
	}

	return issues
}

func isStdlibPackage(path string) bool {
	// Simple heuristic: stdlib packages don't have dots in the first path segment
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return false
	}
	first := parts[0]
	return !strings.Contains(first, ".")
}
