package checks

import (
	"codequality"
	"regexp"
	"strings"
)

type hardcodedSecretsCheck struct{}

// HardcodedSecrets creates a new hardcodedsecrets check
func HardcodedSecrets() codequality.Check {
	return &hardcodedSecretsCheck{}
}

func (c *hardcodedSecretsCheck) Name() string {
	return "hardcoded-secrets"
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*"[^"]+"`),
	regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*"[^"]+"`),
	regexp.MustCompile(`(?i)(secret|token|auth)\s*[:=]\s*"[^"]+"`),
	regexp.MustCompile(`(?i)(access[_-]?key|accesskey)\s*[:=]\s*"[^"]+"`),
	regexp.MustCompile(`(?i)(private[_-]?key|privatekey)\s*[:=]\s*"[^"]+"`),
}

func (c *hardcodedSecretsCheck) Run(ctx *codequality.Context) []codequality.Issue {
	var issues []codequality.Issue

	for _, f := range ctx.Files {
		lines := strings.Split(string(f.Src), "\n")
		for lineNum, line := range lines {
			// Skip comments
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
				continue
			}

			for _, pattern := range secretPatterns {
				if pattern.MatchString(line) {
					issues = append(issues, codequality.Issue{
						Rule:     c.Name(),
						Severity: codequality.SeverityError,
						Message:  "Possible hardcoded secret detected",
						File:     f.Path,
						Line:     lineNum + 1,
						Column:   1,
					})
					break // Only report once per line
				}
			}
		}
	}

	return issues
}
