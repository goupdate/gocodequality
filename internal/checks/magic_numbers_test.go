package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestMagicNumbers(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "magic number in code",
			code: `package main

func main() {
	x := 42
	_ = x
}
`,
			expectIssues: 1,
		},
		{
			name: "zero and one not flagged",
			code: `package main

func main() {
	x := 0
	y := 1
	_ = x
	_ = y
}
`,
			expectIssues: 0,
		},
		{
			name: "constant not flagged",
			code: `package main

const MaxSize = 100

func main() {
	x := MaxSize
	_ = x
}
`,
			expectIssues: 0,
		},
		{
			name: "multiple magic numbers",
			code: `package main

func main() {
	x := 42
	y := 100
	_ = x
	_ = y
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

			check := MagicNumbers()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
