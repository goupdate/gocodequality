package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestImportOrder(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "correct order - stdlib first",
			code: `package main

import (
	"fmt"
	"os"
	"github.com/example/pkg"
)

func main() {
	fmt.Println(os.Args)
	_ = pkg.X
}
`,
			expectIssues: 0,
		},
		{
			name: "incorrect order - external before stdlib",
			code: `package main

import (
	"github.com/example/pkg"
	"fmt"
)

func main() {
	fmt.Println(pkg.X)
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

			check := ImportOrder()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
