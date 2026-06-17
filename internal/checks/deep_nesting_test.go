package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestDeepNesting(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "shallow nesting",
			code: `package main

func foo() {
	if true {
		x := 1
		_ = x
	}
}
`,
			expectIssues: 0,
		},
		{
			name: "deep nesting over 4 levels",
			code: `package main

func foo() {
	if true {
		if true {
			if true {
				if true {
					if true {
						x := 1
						_ = x
					}
				}
			}
		}
	}
}
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

			check := DeepNesting()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
