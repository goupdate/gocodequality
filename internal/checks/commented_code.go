package checks

import (
	"codequality"
	"go/ast"
	"regexp"
	"strings"
)

type commentedCodeCheck struct{}

// CommentedCode creates a new commentedcode check
func CommentedCode() codequality.Check {
	return &commentedCodeCheck{}
}

// Name returns the rule name for this check
func (c *commentedCodeCheck) Name() string {
	return "commented-code"
}

// Patterns that look like Go code in comments
var codePatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\s*(func|var|const|type|import|package)\s+`),
	regexp.MustCompile(`^\s*\w+\s*[:=]\s*`),
	regexp.MustCompile(`^\s*(if|for|switch|return|defer|go)\s+`),
	regexp.MustCompile(`^\s*\w+\.\w+\(`),
	regexp.MustCompile(`^\s*(fmt\.|log\.|os\.|io\.)`),
	regexp.MustCompile(`^\s*\}\s*$`),
	regexp.MustCompile(`^\s*\{\s*$`),
}

// Run executes the check and returns found issues
func (c *commentedCodeCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		// Collect all doc comments to skip them
		docComments := make(map[*ast.CommentGroup]bool)
		for _, decl := range f.AST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Doc != nil {
					docComments[d.Doc] = true
				}
			case *ast.GenDecl:
				if d.Doc != nil {
					docComments[d.Doc] = true
				}
			}
		}

		for _, cg := range f.AST.Comments {
			// Skip doc comments
			if docComments[cg] {
				continue
			}

			for _, comment := range cg.List {
				text := comment.Text
				if !strings.HasPrefix(text, "//") {
					continue
				}
				code := strings.TrimPrefix(text, "//")
				code = strings.TrimSpace(code)

				if code == "" {
					continue
				}

				for _, pattern := range codePatterns {
					if pattern.MatchString(code) {
						pos := f.FSet.Position(comment.Pos())
						issues = append(issues, codequality.Issue{
							Rule:     c.Name(),
							Severity: codequality.SeverityInfo,
							Message:  "Commented-out code detected: " + truncate(code, 60),
							File:     f.Path,
							Line:     pos.Line,
							Column:   pos.Column,
						})
						break
					}
				}
			}
		}
	}

	return issues
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
