package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestLongFunctions(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "short function",
			code: `package main

func short() {
	x := 1
	_ = x
}
`,
			expectIssues: 0,
		},
		{
			name: "long function over 100 lines",
			code: `package main

func longFunc() {
` + strings.Repeat("\tx := 1\n\t_ = x\n", 60) + `}
`,
			expectIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			ctx := &codequality.Context{
				Files: []*codequality.ParsedFile{
					{Path: "test.go", FSet: fset, AST: node, Src: []byte(tt.code)},
				},
			}

			check := LongFunctions()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
