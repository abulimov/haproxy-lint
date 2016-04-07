package checks

import "github.com/abulimov/haproxy-lint/lib"

// Run runs all checks on Sections
func Run(sections []*lib.Section) []lib.Problem {
	var sectionChecks = []lib.SectionCheck{
		CheckUnusedACL,
		CheckUnknownACLs,
		CheckPrecedence,
		CheckDuplicates,
		CheckDeprecations,
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
