// Package codequality provides a pluggable code quality analysis framework for Go projects.
//
// It parses Go source files, builds a shared context, and runs registered checks
// that each return a list of issues found.
package codequality

import (
	"go/ast"
	"go/token"
)

// Severity represents the importance of a found issue.
type Severity int

const (
	SeverityInfo    Severity = iota
	SeverityWarning
	SeverityError
)

// String returns the string representation of the severity level.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "info"
	case SeverityWarning:
		return "warning"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}

// Issue represents a single problem found by a check.
type Issue struct {
	Rule     string
	Severity Severity
	Message  string
	File     string
	Line     int
	Column   int
}

// ParsedFile holds the result of parsing a single Go source file.
type ParsedFile struct {
	Path  string
	FSet  *token.FileSet
	AST   *ast.File
	Src   []byte
}

// Context is passed to each check with all parsed files and project metadata.
type Context struct {
	RootDir string
	Files   []*ParsedFile
}

// Check is the interface every code quality check must implement.
type Check interface {
	Name() string
	Run(ctx *Context) []Issue
}

// Runner orchestrates running multiple checks against a project.
type Runner struct {
	checks []Check
}

// NewRunner creates a Runner with the given checks registered.
func NewRunner(checks ...Check) *Runner {
	return &Runner{checks: checks}
}

// Run executes all registered checks and returns the combined list of issues.
func (r *Runner) Run(ctx *Context) []Issue {
	var issues []Issue
	for _, c := range r.checks {
		issues = append(issues, c.Run(ctx)...)
	}
	return issues
}
