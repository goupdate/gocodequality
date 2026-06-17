package checks

import (
	"codequality"
	"go/ast"
	"go/token"
	"strings"
)

type unusedGlobalVarsCheck struct{}

// UnusedGlobalVars creates a new unusedglobalvars check
func UnusedGlobalVars() codequality.Check {
	return &unusedGlobalVarsCheck{}
}

func (c *unusedGlobalVarsCheck) Name() string {
	return "unused-global-vars"
}

func (c *unusedGlobalVarsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	declarations := make(map[string][]varDecl)
	usages := make(map[string]int)

	for _, f := range ctx.Files {
		collectVarDecls(f.AST, f.FSet, f.Path, declarations)
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
					Message:  "Global variable '" + name + "' is declared but never used",
					File:     decl.path,
					Line:     decl.line,
					Column:   decl.column,
				})
			}
		}
	}

	return issues
}

type varDecl struct {
	path     string
	line     int
	column   int
	isExport bool
}

func collectVarDecls(node *ast.File, fset *token.FileSet, path string, declarations map[string][]varDecl) {
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range vs.Names {
				if strings.HasPrefix(name.Name, "_") {
					continue
				}
				pos := fset.Position(name.Pos())
				declarations[name.Name] = append(declarations[name.Name], varDecl{
					path:     path,
					line:     pos.Line,
					column:   pos.Column,
					isExport: name.IsExported(),
				})
			}
		}
	}
}
