package checks

import (
	"codequality"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnbufferedChannels(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectIssues int
	}{
		{
			name: "unbuffered channel",
			code: `package main

func main() {
	ch := make(chan int)
	_ = ch
}
`,
			expectIssues: 1,
		},
		{
			name: "buffered channel",
			code: `package main

func main() {
	ch := make(chan int, 10)
	_ = ch
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

			check := UnbufferedChannels()
			issues := check.Run(ctx)

			if len(issues) != tt.expectIssues {
				t.Errorf("Expected %d issues, got %d: %+v", tt.expectIssues, len(issues), issues)
			}
		})
	}
}
