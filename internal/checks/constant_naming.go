package checks

import (
	"codequality"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

type constantNamingCheck struct{}

// ConstantNaming creates a new constantnaming check
func ConstantNaming() codequality.Check {
	return &constantNamingCheck{}
}

func (c *constantNamingCheck) Name() string {
	return "constant-naming"
}

var upperSnakeCaseRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

func (c *constantNamingCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}

			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, name := range vs.Names {
					if !name.IsExported() {
						continue
					}

					// Check if name follows UPPER_SNAKE_CASE
					if !upperSnakeCaseRe.MatchString(name.Name) && !isCamelCase(name.Name) {
						pos := f.FSet.Position(name.Pos())
						issues = append(issues, codequality.Issue{
							Rule:     c.Name(),
							Severity: codequality.SeverityInfo,
							Message:  "Exported constant '" + name.Name + "' should use UPPER_SNAKE_CASE or CamelCase",
							File:     f.Path,
							Line:     pos.Line,
							Column:   pos.Column,
						})
					}
				}
			}
		}
	}

	return issues
}

func isCamelCase(s string) bool {
	// Allow standard Go CamelCase (exported identifiers)
	if len(s) == 0 {
		return false
	}
	// Must start with uppercase
	if !isUpper(rune(s[0])) {
		return false
	}
	// Should not contain underscores (except for UPPER_SNAKE_CASE which is checked separately)
	return !strings.Contains(s, "_")
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}
