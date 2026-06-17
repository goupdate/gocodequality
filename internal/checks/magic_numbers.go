package checks

import (
	"codequality"
	"go/ast"
	"go/token"
)

type magicNumbersCheck struct{}

// MagicNumbers creates a new magicnumbers check
func MagicNumbers() codequality.Check {
	return &magicNumbersCheck{}
}

func (c *magicNumbersCheck) Name() string {
	return "magic-numbers"
}

func (c *magicNumbersCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			// Skip const declarations
			if _, ok := n.(*ast.GenDecl); ok {
				if gd, ok := n.(*ast.GenDecl); ok && gd.Tok == token.CONST {
					return false
				}
			}

			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.INT {
				return true
			}

			// Allow 0 and 1 as they are common
			if lit.Value == "0" || lit.Value == "1" {
				return true
			}

			pos := f.FSet.Position(lit.Pos())
			issues = append(issues, codequality.Issue{
				Rule:     c.Name(),
				Severity: codequality.SeverityInfo,
				Message:  "Magic number " + lit.Value + " - consider using a named constant",
				File:     f.Path,
				Line:     pos.Line,
				Column:   pos.Column,
			})

			return true
		})
	}

	return issues
}
