package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestTooManyParams(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "few parameters",
			code: `package main

func foo(a, b int) {}
`,
			expectIssues: 0,
		},
		{
			name: "too many parameters",
			code: `package main

func foo(a, b, c, d, e, f int) {}
`,
			expectIssues: 1,
		},
		{
			name: "exactly 5 parameters",
			code: `package main

func foo(a, b, c, d, e int) {}
`,
			expectIssues: 0,
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

			check := TooManyParams()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
