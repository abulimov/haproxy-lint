package checks

import "github.com/abulimov/haproxy-lint/lib"

func Run(sections []*lib.Section) []lib.Problem {
	var sectionChecks = []lib.SectionCheck{
		CheckUnusedACL,
		CheckUnknownACLs,
		CheckPrecedence,
		CheckDuplicates,
	}
	var globalChecks = []lib.GlobalCheck{
		CheckUnusedBackends,
		CheckUnknownBackends,
	}
	var problems []lib.Problem
	for _, s := range sections {
		for _, r := range sectionChecks {
			problems = append(problems, r(s)...)
		}
	}
	for _, r := range globalChecks {
		problems = append(problems, r(sections)...)
	}
	return problems
}
