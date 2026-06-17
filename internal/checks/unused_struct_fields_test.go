package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnusedStructFields(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "unused private field",
			code: `package main

type MyStruct struct {
	unusedField int
}
`,
			expectIssues: 1,
		},
		{
			name: "used private field",
			code: `package main

type MyStruct struct {
	usedField int
}

func main() {
	s := MyStruct{usedField: 1}
	_ = s.usedField
}
`,
			expectIssues: 0,
		},
		{
			name: "exported field not flagged",
			code: `package main

type MyStruct struct {
	ExportedField int
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

			check := UnusedStructFields()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
