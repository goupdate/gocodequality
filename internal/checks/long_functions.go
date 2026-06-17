package checks

import (
	"codequality"
	"go/ast"
)

const (
	// DefaultMaxFunctionLines is the default maximum number of lines for a function
	DefaultMaxFunctionLines = 100
)

type longFunctionsCheck struct {
	maxLines int
}

// LongFunctions creates a check that flags functions longer than 100 lines
func LongFunctions() codequality.Check {
	return &longFunctionsCheck{maxLines: DefaultMaxFunctionLines}
}

func (c *longFunctionsCheck) Name() string {
	return "long-functions"
}

func (c *longFunctionsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			fd, ok := n.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				return true
			}

			startPos := f.FSet.Position(fd.Body.Lbrace)
			endPos := f.FSet.Position(fd.Body.Rbrace)
			lines := endPos.Line - startPos.Line

			if lines > c.maxLines {
				pos := f.FSet.Position(fd.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityWarning,
					Message:  "Function '" + fd.Name.Name + "' is " + itoa(lines) + " lines (max " + itoa(c.maxLines) + ")",
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}

			return true
		})
	}

	return issues
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
