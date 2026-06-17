package internal

import (
	"codequality"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParseProject walks the given root directory and parses all .go files (excluding tests).
// Returns a Context with all parsed files ready for checks.
func ParseProject(rootDir string) (*codequality.Context, error) {
	var files []*codequality.ParsedFile

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip common non-source directories
		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" || 
			   name == ".plans" || name == ".omo" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files, skip test files
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Read source
		src, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip unreadable files
		}

		// Parse file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			return nil // Skip unparseable files
		}

		files = append(files, &codequality.ParsedFile{
			Path: path,
			FSet: fset,
			AST:  node,
			Src:  src,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &codequality.Context{
		RootDir: rootDir,
		Files:   files,
	}, nil
}
