package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestHardcodedSecrets(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "hardcoded password",
			code: `package main

var password = "secret123"
`,
			expectIssues: 1,
		},
		{
			name: "hardcoded api key",
			code: `package main

var apiKey = "abc123xyz"
`,
			expectIssues: 1,
		},
		{
			name: "no secrets",
			code: `package main

var username = "admin"
`,
			expectIssues: 0,
		},
		{
			name: "secret in comment not flagged",
			code: `package main

// password = "secret"
var x = 1
`,
			expectIssues: 0,
		},
		{
			name: "hardcoded token",
			code: `package main

const token = "mytoken123"
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

			check := HardcodedSecrets()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
