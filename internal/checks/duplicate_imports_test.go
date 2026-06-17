package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestDuplicateImports(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "no duplicates",
			code: `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Args)
}
`,
			expectIssues: 0,
		},
		{
			name: "duplicate import",
			code: `package main

import (
	"fmt"
	"fmt"
)

func main() {
	fmt.Println("hello")
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

			check := DuplicateImports()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
