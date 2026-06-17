package codequality_test

import (
	"codequality"
	"codequality/internal"
	"codequality/internal/checks"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestCodeQualityOnSelf runs all checks on the codequality project itself
func TestCodeQualityOnSelf(t *testing.T) {
	// Find project root
	rootDir := findProjectRoot(t)

	// Parse project
	ctx, err := internal.ParseProject(rootDir)
	if err != nil {
		t.Fatalf("Failed to parse project: %v", err)
	}

	// Register all checks
	allChecks := []codequality.Check{
		checks.UnusedFunctions(),
		checks.EmptyBodyFunctions(),
		checks.UnusedGlobalVars(),
		checks.UnusedStructFields(),
		checks.TodoFixme(),
		checks.CommentedCode(),
		checks.LongFunctions(),
		checks.MissingGodoc(),
		checks.IgnoredErrors(),
		checks.HardcodedSecrets(),
		checks.PanicUsage(),
		checks.MagicNumbers(),
		checks.DeepNesting(),
		checks.TooManyParams(),
		checks.CyclomaticComplexity(),
		checks.UnsafeImport(),
		checks.ObjectInLoop(),
		checks.UnbufferedChannels(),
		checks.ImportOrder(),
		checks.DuplicateImports(),
		checks.ConstantNaming(),
	}

	// Run all checks
	runner := codequality.NewRunner(allChecks...)
	issues := runner.Run(ctx)

	// Report all issues
	if len(issues) > 0 {
		t.Logf("Found %d issues in codequality project:", len(issues))
		for _, issue := range issues {
			relPath, _ := filepath.Rel(rootDir, issue.File)
			t.Logf("  %s:%d:%d: [%s] %s: %s",
				relPath, issue.Line, issue.Column,
				issue.Severity, issue.Rule, issue.Message)
		}
	}

	// Group by rule
	ruleCounts := make(map[string]int)
	for _, issue := range issues {
		ruleCounts[issue.Rule]++
	}

	t.Logf("\nIssues by rule:")
	for rule, count := range ruleCounts {
		t.Logf("  %s: %d", rule, count)
	}

	// Fail if there are errors
	errorCount := 0
	for _, issue := range issues {
		if issue.Severity == codequality.SeverityError {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("Found %d errors in codequality project", errorCount)
	}
}

func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root (go.mod)")
		}
		dir = parent
	}
}

// Benchmark to show performance
func BenchmarkCodeQuality(b *testing.B) {
	rootDir := findProjectRootFromT(b)
	ctx, err := internal.ParseProject(rootDir)
	if err != nil {
		b.Fatalf("Failed to parse project: %v", err)
	}

	allChecks := []codequality.Check{
		checks.UnusedFunctions(),
		checks.EmptyBodyFunctions(),
		checks.UnusedGlobalVars(),
		checks.UnusedStructFields(),
		checks.TodoFixme(),
		checks.CommentedCode(),
		checks.LongFunctions(),
		checks.MissingGodoc(),
		checks.IgnoredErrors(),
		checks.HardcodedSecrets(),
		checks.PanicUsage(),
		checks.MagicNumbers(),
		checks.DeepNesting(),
		checks.TooManyParams(),
		checks.CyclomaticComplexity(),
		checks.UnsafeImport(),
		checks.ObjectInLoop(),
		checks.UnbufferedChannels(),
		checks.ImportOrder(),
		checks.DuplicateImports(),
		checks.ConstantNaming(),
	}

	runner := codequality.NewRunner(allChecks...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = runner.Run(ctx)
	}
}

func findProjectRootFromT(b *testing.B) string {
	dir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			b.Fatal("Could not find project root (go.mod)")
		}
		dir = parent
	}
}

// Helper to print issues in a readable format
func printIssues(issues []codequality.Issue, rootDir string) {
	fmt.Printf("\n=== Code Quality Issues ===\n\n")
	for _, issue := range issues {
		relPath, _ := filepath.Rel(rootDir, issue.File)
		fmt.Printf("%s:%d:%d: [%s] %s\n",
			relPath, issue.Line, issue.Column,
			issue.Severity, issue.Rule)
		fmt.Printf("  %s\n\n", issue.Message)
	}
}
