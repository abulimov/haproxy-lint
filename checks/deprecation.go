package checks

import (
	"fmt"

	"github.com/abulimov/haproxy-lint/lib"
)

// CheckDeprecations checks if we have duplicated rules
func CheckDeprecations(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	deprecations := map[string]string{
		"block":              "http-request deny",
		"clitimeout":         "timeout client",
		"timeout clitimeout": "timeout client",
		"contimeout":         "timeout connect",
		"timeout contimeout": "timeout connect",
		"srvtimeout":         "timeout server",
		"timeout srvtimeout": "timeout server",
		"redisp":             "option redispatch",
		"redispatch":         "option redispatch",
		"transparent":        "option transparent",
	}
	for i, line := range s.Content {
		for deprecated, replacement := range deprecations {
			found := lib.GetUsage(deprecated, line)
			if found != "" {
				problems = append(
					problems,
					lib.Problem{
						Line:     i,
						Col:      0,
						Severity: "warning",
						Message:  fmt.Sprintf("The '%s' directive is now deprecated in favor of '%s'", deprecated, replacement),
					},
				)
			}
		}
	}
	return problems
}
