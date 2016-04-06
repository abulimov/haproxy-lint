package checks

import (
	"fmt"
	"sort"

	"github.com/abulimov/haproxy-lint/lib"
)

// CheckDuplicates checks if we have duplicated rules
func CheckDuplicates(s *lib.Section) []lib.Problem {
	var problems []lib.Problem
	seen := map[string]int{}
	var lineNums []int
	for i := range s.Content {
		lineNums = append(lineNums, i)
	}
	sort.Ints(lineNums)
	for _, i := range lineNums {
		trimmed := lib.StripComments(s.Content[i])
		if trimmed != "" {
			prev, found := seen[trimmed]
			if found {
				problems = append(
					problems,
					lib.Problem{
						Line:     i,
						Col:      0,
						Severity: "warning",
						Message:  fmt.Sprintf("a '%s' rule is duplicate of line %d", trimmed, prev),
					},
				)
			} else {
				seen[trimmed] = i
			}
		}
	}
	return problems
}
