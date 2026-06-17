# Code Quality

A fast, extensible static code analysis tool for Go projects. Analyzes your codebase for common issues, code smells, and style violations.

## Features

### Code Quality Checks (21 rules)

**High Priority**
- `todo-fixme` — Detects TODO, FIXME, HACK, XXX comments
- `commented-code` — Finds commented-out code blocks
- `long-functions` — Flags functions longer than 100 lines
- `missing-godoc` — Checks for missing documentation on exported functions and types

**Medium Priority**
- `ignored-errors` — Detects `_ = fn()` patterns where errors are ignored
- `hardcoded-secrets` — Finds hardcoded passwords, API keys, tokens
- `panic-usage` — Flags use of `panic()` outside initialization
- `magic-numbers` — Detects numeric literals without named constants

**Low Priority**
- `deep-nesting` — Flags functions with nesting depth > 4
- `too-many-params` — Flags functions with more than 5 parameters
- `cyclomatic-complexity` — Detects functions with complexity > 10
- `unsafe-import` — Flags imports of the `unsafe` package
- `object-in-loop` — Detects allocations (`make`, `new`, `time.Now`) inside loops
- `unbuffered-channels` — Flags unbuffered channel creation
- `import-order` — Checks stdlib imports come before external
- `duplicate-imports` — Detects duplicate import statements
- `constant-naming` — Validates exported constants use UPPER_SNAKE_CASE or CamelCase

**Original Checks**
- `unused-functions` — Detects declared but never used private functions
- `empty-body-functions` — Flags functions with empty bodies
- `unused-global-vars` — Detects unused global variables
- `unused-struct-fields` — Finds unused struct fields

## Installation

```bash
go install github.com/yourusername/codequality/cmd/codequality@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/codequality.git
cd codequality
go build -o codequality ./cmd/codequality
```

## Usage

```bash
# Analyze a project
codequality ./path/to/project

# Analyze current directory
codequality .
```

### Output Format

```
file.go:10:5: [error] hardcoded-secrets: Possible hardcoded secret detected
file.go:20:1: [warning] long-functions: Function 'Handle' is 150 lines (max 100)
file.go:30:2: [info] todo-fixme: TODO comment found: TODO: implement this

Found 3 issues (1 errors, 1 warnings, 1 info)
```

Exit codes:
- `0` — No errors found (warnings/info only or clean)
- `1` — Errors found
- `2` — Invalid arguments or parse error

## Adding a New Check

1. Create `internal/checks/my_check.go`:

```go
package checks

import "codequality"

type myCheck struct{}

// MyCheck creates a new check for ...
func MyCheck() codequality.Check {
    return &myCheck{}
}

// Name returns the rule name
func (c *myCheck) Name() string {
    return "my-check"
}

// Run executes the check
func (c *myCheck) Run(ctx *codequality.Context) []codequality.Issue {
    var issues []codequality.Issue
    // Your analysis logic here
    return issues
}
```

2. Create tests in `internal/checks/my_check_test.go`

3. Register in `cmd/codequality/main.go`:

```go
allChecks := []codequality.Check{
    // ... existing checks
    checks.MyCheck(),
}
```

## Library Usage

Use as a library in your own tools:

```go
package main

import (
    "codequality"
    "codequality/internal"
    "codequality/internal/checks"
    "fmt"
)

func main() {
    ctx, err := internal.ParseProject("./myproject")
    if err != nil {
        panic(err)
    }

    runner := codequality.NewRunner(
        checks.UnusedFunctions(),
        checks.LongFunctions(),
        checks.HardcodedSecrets(),
    )

    issues := runner.Run(ctx)
    for _, issue := range issues {
        fmt.Printf("%s:%d: [%s] %s: %s\n",
            issue.File, issue.Line,
            issue.Severity, issue.Rule, issue.Message)
    }
}
```

## Testing

```bash
# Run all tests
go test ./...

# Run self-analysis test
go test -v -run TestCodeQualityOnSelf
```

## Project Structure

```
codequality/
├── codequality.go          # Public API: Runner, Check, Issue, Context
├── go.mod
├── internal/
│   ├── parser.go           # Walk + parse with caching
│   └── checks/             # 21 checks + tests
│       ├── unused_functions.go
│       ├── unused_functions_test.go
│       └── ...
├── cmd/
│   └── codequality/
│       └── main.go         # CLI entry point
└── testdata/               # Test fixtures
```

## License

MIT
