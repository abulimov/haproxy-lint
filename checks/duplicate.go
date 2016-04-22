package checks

import (
	"fmt"
	"sort"
	"strings"

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
		line := s.Content[i]
		normalized := strings.Join(strings.Fields(line), " ")
		if normalized != "" {
			prev, found := seen[normalized]
			if found {
				problems = append(
					problems,
					lib.Problem{
						Line:     i,
						Col:      0,
						Severity: "warning",
						Message:  fmt.Sprintf("a '%s' rule is duplicate of line %d", line, prev),
					},
				)
			} else {
				seen[line] = i
			}
		}
	}
	return problems
}
