package checks

import (
	"codequality"
	"go/ast"
	"strings"
)

type missingGodocCheck struct{}

// MissingGodoc creates a new missinggodoc check
func MissingGodoc() codequality.Check {
	return &missingGodocCheck{}
}

func (c *missingGodocCheck) Name() string {
	return "missing-godoc"
}

func (c *missingGodocCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Only check exported functions
			if !fd.Name.IsExported() {
				continue
			}

			// Skip test/benchmark/example functions
			name := fd.Name.Name
			if strings.HasPrefix(name, "Test") ||
				strings.HasPrefix(name, "Benchmark") ||
				strings.HasPrefix(name, "Example") {
				continue
			}

			// Check if function has a doc comment
			if fd.Doc == nil || len(fd.Doc.List) == 0 {
				pos := f.FSet.Position(fd.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityInfo,
					Message:  "Exported function '" + name + "' should have a godoc comment",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}
		}

		// Check exported types
		for _, decl := range f.AST.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !ts.Name.IsExported() {
					continue
				}

				doc := ts.Doc
				if doc == nil {
					doc = gd.Doc
				}
				if doc == nil || len(doc.List) == 0 {
					pos := f.FSet.Position(ts.Pos())
					issues = append(issues, codequality.Issue{
						Rule:     c.Name(),
						Severity: codequality.SeverityInfo,
						Message:  "Exported type '" + ts.Name.Name + "' should have a godoc comment",
						File:     f.Path,
						Line:     pos.Line,
						Column:   pos.Column,
					})
				}
			}
		}
	}

	return issues
}
