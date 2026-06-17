package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestPanicUsage(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "panic usage",
			code: `package main

func main() {
	panic("something went wrong")
}
`,
			expectIssues: 1,
		},
		{
			name: "no panic",
			code: `package main

func main() {
	x := 1
	_ = x
}
`,
			expectIssues: 0,
		},
		{
			name: "multiple panics",
			code: `package main

func foo() {
	if true {
		panic("error 1")
	}
	panic("error 2")
}
`,
			expectIssues: 2,
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

			check := PanicUsage()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
