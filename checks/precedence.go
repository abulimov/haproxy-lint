package checks

import (
	"fmt"
	"sort"

	"github.com/abulimov/haproxy-lint/lib"
)

// CheckPrecedence checks if we have wrong ordered rules in our section
func CheckPrecedence(s *lib.Section) []lib.Problem {
	keywordsMap := map[string]int{
		"http-request": 1,
		"reqrep":       1,
		"reqirep":      1,
		"redirect":     2,
		"use-server":   3,
		"use_backend":  3,
	}
	var problems []lib.Problem
	var lineNums []int
	keywords := make(map[int]string)
	for i, line := range s.Content {
		for kw := range keywordsMap {
			name := lib.GetUsage(kw, line)
			if name != "" {
				lineNums = append(lineNums, i)
				keywords[i] = kw
				break
			}
		}
	}
	sort.Ints(lineNums)
	maxPriority := 0
	var maxKW string
	for _, i := range lineNums {
		kw := keywords[i]
		curPriority := keywordsMap[kw]
		if curPriority >= maxPriority {
			maxPriority = curPriority
			maxKW = kw
		} else {
			problems = append(
				problems,
				lib.Problem{
					Line:     i,
					Col:      0,
					Severity: "warning",
					Message:  fmt.Sprintf("a '%s' rule placed after a '%s' rule will still be processed before", kw, maxKW),
				},
			)
		}
	}
	return problems
}
