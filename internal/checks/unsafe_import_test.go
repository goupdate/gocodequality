package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnsafeImport(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "unsafe import",
			code: `package main

import "unsafe"

func main() {
	_ = unsafe.Sizeof(0)
}
`,
			expectIssues: 1,
		},
		{
			name: "no unsafe import",
			code: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
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

			check := UnsafeImport()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
