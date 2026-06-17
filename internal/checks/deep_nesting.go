package checks

import (
	"codequality"
	"go/ast"
)

const (
	// DefaultMaxNestingDepth is the default maximum nesting depth
	DefaultMaxNestingDepth = 4
)

type deepNestingCheck struct {
	maxDepth int
}

// DeepNesting creates a check that flags functions with nesting depth > 4
func DeepNesting() codequality.Check {
	return &deepNestingCheck{maxDepth: DefaultMaxNestingDepth}
}

func (c *deepNestingCheck) Name() string {
	return "deep-nesting"
}

func (c *deepNestingCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		for _, decl := range f.AST.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Body == nil {
				continue
			}

			c.checkNesting(fd.Body, 0, fd, f, &issues)
		}
	}

	return issues
}

func (c *deepNestingCheck) checkNesting(node ast.Node, depth int, fd *ast.FuncDecl, f *codequality.ParsedFile, issues *[]codequality.Issue) {
	switch n := node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.SelectStmt, *ast.TypeSwitchStmt:
		depth++
		if depth > c.maxDepth {
			pos := f.FSet.Position(n.Pos())
			*issues = append(*issues, codequality.Issue{
				Rule:     c.Name(),
				Severity: codequality.SeverityWarning,
				Message:  "Function '" + fd.Name.Name + "' has nesting depth " + itoa(depth) + " (max " + itoa(c.maxDepth) + ")",
				File:     f.Path,
				Line:     pos.Line,
				Column:   pos.Column,
			})
		}
	}

	ast.Inspect(node, func(child ast.Node) bool {
		if child == node {
			return true
		}
		switch child.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.SwitchStmt, *ast.SelectStmt, *ast.TypeSwitchStmt:
			c.checkNesting(child, depth, fd, f, issues)
			return false
		}
		return true
	})
}
