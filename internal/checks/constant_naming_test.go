package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestConstantNaming(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "UPPER_SNAKE_CASE constant",
			code: `package main

const MAX_SIZE = 100
`,
			expectIssues: 0,
		},
		{
			name: "CamelCase constant",
			code: `package main

const MaxSize = 100
`,
			expectIssues: 0,
		},
		{
			name: "private constant not checked",
			code: `package main

const maxSize = 100
`,
			expectIssues: 0,
		},
		{
			name: "mixed case with underscore",
			code: `package main

const Max_Size = 100
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

			check := ConstantNaming()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
