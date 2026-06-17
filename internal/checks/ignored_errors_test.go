package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestIgnoredErrors(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "ignored error with blank identifier",
			code: `package main

import "os"

func main() {
	_ = os.Remove("file.txt")
}
`,
			expectIssues: 1,
		},
		{
			name: "error handled properly",
			code: `package main

import "os"

func main() {
	err := os.Remove("file.txt")
	if err != nil {
		panic(err)
	}
}
`,
			expectIssues: 0,
		},
		{
			name: "multiple return values handled",
			code: `package main

func main() {
	result, err := someFunc()
	_ = result
	if err != nil {
		panic(err)
	}
}

func someFunc() (int, error) {
	return 0, nil
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

			check := IgnoredErrors()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
