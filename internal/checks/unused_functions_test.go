package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnusedFunctions(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectIssues  int
	}{
		{
			name: "unused private function",
			code: `package main

func unusedFunc() {}
`,
			expectIssues: 1,
		},
		{
			name: "used private function",
			code: `package main

func usedFunc() {}

func main() {
	usedFunc()
}
`,
			expectIssues: 0,
		},
		{
			name: "exported function not flagged",
			code: `package main

func ExportedFunc() {}
`,
			expectIssues: 0,
		},
		{
			name: "test function not flagged",
			code: `package main

func TestSomething() {}
`,
			expectIssues: 0,
		},
		{
			name: "main function not flagged",
			code: `package main

func main() {}
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

			check := UnusedFunctions()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
