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
	// channel to get problems from check
	ch := make(chan []lib.Problem)
	defer close(ch)
	// count how many times we call goroutines
	count := 0
	for _, s := range sections {
		for _, r := range sectionChecks {
			if !extrasOnly || r.extra {
				count++
				go func(s *lib.Section, r secCheck) {
					ch <- r.f(s)
				}(s, r)
			}
		}
	}
	for _, r := range globalChecks {
		if !extrasOnly || r.extra {
			count++
			go func(r globCheck) {
				ch <- r.f(sections)
			}(r)
		}
	}
	for i := 0; i < count; i++ {
		p := <-ch
		problems = append(problems, p...)
	}
	return problems
}
