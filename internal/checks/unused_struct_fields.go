package checks

import (
	"codequality"
	"go/ast"
)

type unusedStructFieldsCheck struct{}

// UnusedStructFields creates a new unusedstructfields check
func UnusedStructFields() codequality.Check {
	return &unusedStructFieldsCheck{}
}

func (c *unusedStructFieldsCheck) Name() string {
	return "unused-struct-fields"
}

func (c *unusedStructFieldsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	type fieldDecl struct {
		path     string
		line     int
		column   int
		isExport bool
		fullName string
	}

	declarations := make(map[string][]fieldDecl)
	usages := make(map[string]int)

	for _, f := range ctx.Files {
		// Collect struct field declarations
		ast.Inspect(f.AST, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			for _, field := range st.Fields.List {
				for _, name := range field.Names {
					pos := f.FSet.Position(name.Pos())
					fullName := ts.Name.Name + "." + name.Name
					declarations[name.Name] = append(declarations[name.Name], fieldDecl{
						path:     f.Path,
						line:     pos.Line,
						column:   pos.Column,
						isExport: name.IsExported(),
						fullName: fullName,
					})
				}
			}
			return true
		})

		// Collect all usages
		collectAllUsages(f.AST, usages)
	}

	var issues []codequality.Issue

	for name, decls := range declarations {
		for _, decl := range decls {
			if decl.isExport {
				continue
			}
			if usages[name] <= 1 {
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Struct field '" + decl.fullName + "' is declared but never used",
					File:     decl.path,
					Line:     decl.line,
					Column:   decl.column,
				})
			}
		}
	}

	return issues
}
