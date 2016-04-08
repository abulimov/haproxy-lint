package checks

import "github.com/abulimov/haproxy-lint/lib"

// Run runs all checks on Sections
func Run(sections []*lib.Section, extrasOnly bool) []lib.Problem {
	type secCheck struct {
		f lib.SectionCheck
		// extra indicates that this check is not implemented in HAProxy binary
		extra bool
	}
	type globCheck struct {
		f lib.GlobalCheck
		// extra indicates that this check is not implemented in HAProxy binary
		extra bool
	}
	var sectionChecks = []secCheck{
		{CheckUnusedACL, true},
		{CheckUnknownACLs, false},
		{CheckPrecedence, false},
		{CheckDuplicates, true},
		{CheckDeprecations, false},
	}
	var globalChecks = []globCheck{
		{CheckUnusedBackends, true},
		// while check for unknown backend is implemented in HAProxy,
		// alert message for it doesn't show the line with backend usage
		{CheckUnknownBackends, true},
	}
	var problems []lib.Problem
	for _, s := range sections {
		for _, r := range sectionChecks {
			if !extrasOnly || r.extra {
				problems = append(problems, r.f(s)...)
			}
		}
	}
	for _, r := range globalChecks {
		if !extrasOnly || r.extra {
			problems = append(problems, r.f(sections)...)
		}
	}
	return problems
}
