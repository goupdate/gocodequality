package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnusedGlobalVars(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "unused private var",
			code: `package main

var unusedVar = 1
`,
			expectIssues: 1,
		},
		{
			name: "used private var",
			code: `package main

var usedVar = 1

func main() {
	_ = usedVar
}
`,
			expectIssues: 0,
		},
		{
			name: "exported var not flagged",
			code: `package main

var ExportedVar = 1
`,
			expectIssues: 0,
		},
		{
			name: "blank identifier not flagged",
			code: `package main

var _ = 1
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

			check := UnusedGlobalVars()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
