package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestEmptyBodyFunctions(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "function with empty body",
			code: `package main

func emptyFunc() {}
`,
			expectIssues: 1,
		},
		{
			name: "function with body",
			code: `package main

func nonEmptyFunc() {
	x := 1
	_ = x
}
`,
			expectIssues: 0,
		},
		{
			name: "test function with empty body not flagged",
			code: `package main

func TestSomething() {}
`,
			expectIssues: 0,
		},
		{
			name: "multiple functions mixed",
			code: `package main

func empty1() {}

func nonEmpty() {
	x := 1
	_ = x
}

func empty2() {}
`,
			expectIssues: 2,
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

			check := EmptyBodyFunctions()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
