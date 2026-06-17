package checks

import (
	"codequality"
	"go/ast"
	"regexp"
	"strings"
)

type todoFixmeCheck struct{}

// TodoFixme creates a new todofixme check
func TodoFixme() codequality.Check {
	return &todoFixmeCheck{}
}

func (c *todoFixmeCheck) Name() string {
	return "todo-fixme"
}

var todoRe = regexp.MustCompile(`(?i)\b(TODO|FIXME|HACK|XXX)\b`)

func (c *todoFixmeCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, cg := range f.AST.Comments {
			for _, comment := range cg.List {
				text := strings.TrimPrefix(comment.Text, "//")
				text = strings.TrimPrefix(text, "/*")
				text = strings.TrimSuffix(text, "*/")

				loc := todoRe.FindStringIndex(text)
				if loc == nil {
					continue
				}

				keyword := strings.ToUpper(strings.TrimSpace(text[loc[0]:loc[1]]))
				pos := f.FSet.Position(comment.Pos())
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityInfo,
					Message:  keyword + " comment found: " + strings.TrimSpace(text),
					File:     f.Path,
					Line:     pos.Line,
					Column:   pos.Column,
				})
			}
		}

		// Also check doc comments attached to declarations
		ast.Inspect(f.AST, func(n ast.Node) bool {
			return true
		})
	}

	return issues
}
