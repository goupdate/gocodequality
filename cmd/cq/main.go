package main

import (
	"codequality"
	"codequality/internal"
	"codequality/internal/checks"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	dir := os.Args[1]

	// Resolve to absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(2)
	}

	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a valid directory\n", dir)
		os.Exit(2)
	}

	// Parse project
	ctx, err := internal.ParseProject(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing project: %v\n", err)
		os.Exit(2)
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

	// Sort issues by file, then line
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].File != issues[j].File {
			return issues[i].File < issues[j].File
		}
		return issues[i].Line < issues[j].Line
	})

	// Print results
	hasErrors := false
	for _, issue := range issues {
		// Make path relative for cleaner output
		relPath, err := filepath.Rel(absDir, issue.File)
		if err != nil {
			relPath = issue.File
		}
		fmt.Printf("%s:%d:%d: [%s] %s: %s\n",
			relPath, issue.Line, issue.Column,
			issue.Severity, issue.Rule, issue.Message)
		if issue.Severity == codequality.SeverityError {
			hasErrors = true
		}
	}

	// Summary
	errCount, warnCount, infoCount := 0, 0, 0
	for _, issue := range issues {
		switch issue.Severity {
		case codequality.SeverityError:
			errCount++
		case codequality.SeverityWarning:
			warnCount++
		case codequality.SeverityInfo:
			infoCount++
		}
	}

	fmt.Printf("\nFound %d issues (%d errors, %d warnings, %d info)\n",
		len(issues), errCount, warnCount, infoCount)

	if hasErrors {
		os.Exit(1)
	}
}
