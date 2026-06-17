package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestCommentedCode(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "commented function call",
			code: `package main

func foo() {
	// fmt.Println("hello")
}
`,
			expectIssues: 1,
		},
		{
			name: "commented variable assignment",
			code: `package main

func foo() {
	// x := 1
}
`,
			expectIssues: 1,
		},
		{
			name: "normal comment not flagged",
			code: `package main

// This function does something important
func foo() {}
`,
			expectIssues: 0,
		},
		{
			name: "doc comment not flagged",
			code: `package main

// Foo does something
func Foo() {}
`,
			expectIssues: 0,
		},
		{
			name: "commented if statement",
			code: `package main

func foo() {
	// if x > 0 {
}
`,
			expectIssues: 1,
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

			check := CommentedCode()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
