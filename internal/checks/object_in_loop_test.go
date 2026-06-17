package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestObjectInLoop(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "make inside loop",
			code: `package main

func main() {
	for i := 0; i < 10; i++ {
		s := make([]int, 10)
		_ = s
	}
}
`,
			expectIssues: 1,
		},
		{
			name: "make outside loop",
			code: `package main

func main() {
	s := make([]int, 10)
	for i := 0; i < 10; i++ {
		_ = s
	}
}
`,
			expectIssues: 0,
		},
		{
			name: "new inside loop",
			code: `package main

func main() {
	for i := 0; i < 10; i++ {
		p := new(int)
		_ = p
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

			check := ObjectInLoop()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
