package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestMissingGodoc(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "exported function without godoc",
			code: `package main

func ExportedFunc() {}
`,
			expectIssues: 1,
		},
		{
			name: "exported function with godoc",
			code: `package main

// ExportedFunc does something
func ExportedFunc() {}
`,
			expectIssues: 0,
		},
		{
			name: "private function without godoc not flagged",
			code: `package main

func privateFunc() {}
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
			name: "exported type without godoc",
			code: `package main

type ExportedType struct{}
`,
			expectIssues: 1,
		},
		{
			name: "exported type with godoc",
			code: `package main

// ExportedType is something
type ExportedType struct{}
`,
			expectIssues: 0,
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

			check := MissingGodoc()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
