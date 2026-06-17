package checks

import (
	"codequality"
	"go/ast"
	"go/token"
	"strings"
)

type unusedFunctionsCheck struct{}

// UnusedFunctions creates a new unusedfunctions check
func UnusedFunctions() codequality.Check {
	return &unusedFunctionsCheck{}
}

func (c *unusedFunctionsCheck) Name() string {
	return "unused-functions"
}

func (c *unusedFunctionsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	declarations := make(map[string][]funcDecl)
	usages := make(map[string]int)

	for _, f := range ctx.Files {
		collectFuncDecls(f.AST, f.FSet, f.Path, declarations)
		collectAllUsages(f.AST, usages)
	}

	var issues []codequality.Issue

	for name, decls := range declarations {
		for _, decl := range decls {
			if decl.isExport {
				continue
			}
			// usages[name] counts all identifier occurrences including the declaration itself.
			// If count <= 1, the function is only referenced at its own declaration.
			if usages[name] <= 1 {
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Function '" + name + "' is declared but never used",
					File:     decl.path,
					Line:     decl.line,
					Column:   decl.column,
				})
			}
		}
	}

	return issues
}

type funcDecl struct {
	path     string
	line     int
	column   int
	isExport bool
}

func collectFuncDecls(node *ast.File, fset *token.FileSet, path string, declarations map[string][]funcDecl) {
	for _, decl := range node.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := fd.Name.Name

		// Skip test/benchmark/example functions
		if strings.HasPrefix(name, "Test") ||
			strings.HasPrefix(name, "Benchmark") ||
			strings.HasPrefix(name, "Example") {
			continue
		}

		// Skip main
		if name == "main" {
			continue
		}

		pos := fset.Position(fd.Pos())
		declarations[name] = append(declarations[name], funcDecl{
			path:     path,
			line:     pos.Line,
			column:   pos.Column,
			isExport: fd.Name.IsExported(),
		})
	}
}

func collectAllUsages(node *ast.File, usages map[string]int) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			usages[x.Name]++
		case *ast.SelectorExpr:
			usages[x.Sel.Name]++
		}
		return true
	})
}
