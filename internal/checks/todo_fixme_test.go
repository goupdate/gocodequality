package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestTodoFixme(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "TODO comment",
			code: `package main

// TODO: implement this
func foo() {}
`,
			expectIssues: 1,
		},
		{
			name: "FIXME comment",
			code: `package main

// FIXME: broken logic
func foo() {}
`,
			expectIssues: 1,
		},
		{
			name: "no todo comments",
			code: `package main

// This is a normal comment
func foo() {}
`,
			expectIssues: 0,
		},
		{
			name: "multiple TODO comments",
			code: `package main

// TODO: first thing
// TODO: second thing
func foo() {}
`,
			expectIssues: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			ctx := &codequality.Context{
				Files: []*codequality.ParsedFile{
					{Path: "test.go", FSet: fset, AST: node, Src: []byte(tt.code)},
				},
			}

			check := TodoFixme()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
