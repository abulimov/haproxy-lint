package checks

import (
	"fmt"
	"sort"

	"github.com/abulimov/haproxy-lint/lib"
)

// CheckPrecedence checks if we have wrong ordered rules in our section.
// Based on https://raw.githubusercontent.com/tmm1/haproxy-dev/master/src/cfgparse.c
// rules warnif_misplaced_reqadd/warnif_misplaced_reqxxx/warnif_misplaced_block.
func CheckPrecedence(s *lib.Section) []lib.Problem {
	keywordsMap := map[string]int{
		"http-request": 1,
		"reqrep":       2,
		"reqirep":      2,
		"reqdel":       2,
		"reqidel":      2,
		"reqdeny":      2,
		"reqideny":     2,
		"reqpass":      2,
		"reqipass":     2,
		"reqallow":     2,
		"reqiallow":    2,
		"reqtarpit":    2,
		"reqitarpit":   2,
		"reqsetbe":     2,
		"reqisetbe":    2,
		"reqadd":       3,
		"redirect":     4,
		"use-server":   5,
		"use_backend":  5,
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
