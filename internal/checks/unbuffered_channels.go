package checks

import (
	"codequality"
	"go/ast"
)

type unbufferedChannelsCheck struct{}

// UnbufferedChannels creates a new unbufferedchannels check
func UnbufferedChannels() codequality.Check {
	return &unbufferedChannelsCheck{}
}

func (c *unbufferedChannelsCheck) Name() string {
	return "unbuffered-channels"
}

func (c *unbufferedChannelsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		ast.Inspect(f.AST, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if it's make(chan ...)
			ident, ok := call.Fun.(*ast.Ident)
			if !ok || ident.Name != "make" {
				return true
			}

			if len(call.Args) < 1 {
				return true
			}

			// Check if first arg is chan type
			chanType, ok := call.Args[0].(*ast.ChanType)
			if !ok {
				return true
			}

			// Check if buffer size is specified (second arg)
			if len(call.Args) < 2 {
				pos := f.FSet.Position(call.Pos())
				_ = chanType // suppress unused warning
				issues = append(issues, codequality.Issue{
					Rule:     c.Name(),
					Severity: codequality.SeverityInfo,
					Message:  "Unbuffered channel created - may cause blocking",
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
